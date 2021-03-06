package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	errZeroRate     = errors.New("rate must be bigger than zero")
	errEmptyTargets = errors.New("targets is required")
	errBadCert      = errors.New("bad certificate")
)

type csl []string

func (l *csl) UnmarshalJSON(data []byte) (err error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*l = strings.Split(value, ",")
	return nil
}

type headers struct{ http.Header }

func (h headers) UnmarshalJSON(data []byte) (err error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("header '%s' has a wrong format", value)
	}

	key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if key == "" || val == "" {
		return fmt.Errorf("header '%s' has a wrong format", value)
	}
	// Add key/value directly to the http.Header (map[string][]string).
	// http.Header.Add() canonicalizes keys but vegeta is used
	// to test systems that require case-sensitive headers.
	h.Header[key] = append(h.Header[key], val)

	return nil
}

type localAddr struct{ *net.IPAddr }

func (ip *localAddr) UnmarshalJSON(data []byte) (err error) {
	var value string
	if e := json.Unmarshal(data, &value); e != nil {
		return e
	}

	ip.IPAddr, err = net.ResolveIPAddr("ip", value)
	return err
}

// AttackOptions aggregates the vegeta attack options
type AttackOptions struct {
	Targets string
	// Output			string
	Body        string
	Cert        string
	Key         string
	RootCerts   csl
	HTTP2       bool
	Insecure    bool
	Lazy        bool
	Duration    string
	Timeout     string
	Rate        uint64
	Workers     uint64
	Connections int
	Redirects   int
	Headers     headers
	Laddr       localAddr
	Keepalive   bool
}

// NewAttackOptions returns a new AttackOptions with default options
func NewAttackOptions() *AttackOptions {
	return &AttackOptions{
		HTTP2:       true,
		Insecure:    false,
		Lazy:        false,
		Duration:    "0",
		Timeout:     vegeta.DefaultTimeout.String(),
		Rate:        50,
		Workers:     vegeta.DefaultWorkers,
		Connections: vegeta.DefaultConnections,
		Redirects:   vegeta.DefaultRedirects,
		Headers:     headers{http.Header{}},
		Laddr:       localAddr{&vegeta.DefaultLocalAddr},
		Keepalive:   true,
	}
}

// GetAttackWorker generates AttackWorker from AttackOptions
func (opts *AttackOptions) GetAttackWorker(broadcaster Broadcaster) (*AttackWorker, error) {
	if opts.Rate == 0 {
		return nil, errZeroRate
	}

	if opts.Targets == "" {
		return nil, errEmptyTargets
	}

	duration, err := time.ParseDuration(opts.Duration)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(opts.Timeout)
	if err != nil {
		return nil, err
	}

	var (
		tr   vegeta.Targeter
		src  = strings.NewReader(opts.Targets)
		body = []byte(opts.Body)
		hdr  = opts.Headers.Header
	)

	if opts.Lazy {
		tr = vegeta.NewLazyTargeter(src, body, hdr)
	} else if tr, err = vegeta.NewEagerTargeter(src, body, hdr); err != nil {
		return nil, err
	}

	tlsc, err := tlsConfig(opts.Insecure, opts.Cert, opts.Key, opts.RootCerts)
	if err != nil {
		return nil, err
	}

	atk := vegeta.NewAttacker(
		vegeta.Redirects(opts.Redirects),
		vegeta.Timeout(timeout),
		vegeta.LocalAddr(*opts.Laddr.IPAddr),
		vegeta.TLSConfig(tlsc),
		vegeta.Workers(opts.Workers),
		vegeta.KeepAlive(opts.Keepalive),
		vegeta.Connections(opts.Connections),
		vegeta.HTTP2(opts.HTTP2),
	)

	return NewAttackWorker(atk, tr, opts.Rate, duration, broadcaster), nil
}

// tlsConfig builds a *tls.Config from the given options.
// * copied from https://github.com/tsenart/vegeta/blob/master/attack.go
func tlsConfig(insecure bool, certf, keyf string, rootCerts []string) (*tls.Config, error) {
	var err error
	files := map[string][]byte{}
	filenames := append([]string{certf, keyf}, rootCerts...)
	for _, f := range filenames {
		if f != "" {
			if files[f], err = ioutil.ReadFile(f); err != nil {
				return nil, err
			}
		}
	}

	c := tls.Config{InsecureSkipVerify: insecure}
	if cert, ok := files[certf]; ok {
		key, ok := files[keyf]
		if !ok {
			key = cert
		}

		certificate, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		c.Certificates = append(c.Certificates, certificate)
		c.BuildNameToCertificate()
	}

	if len(rootCerts) > 0 {
		c.RootCAs = x509.NewCertPool()
		for _, f := range rootCerts {
			if !c.RootCAs.AppendCertsFromPEM(files[f]) {
				return nil, errBadCert
			}
		}
	}

	return &c, nil
}
