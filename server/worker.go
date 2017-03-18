package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	currentWorker   *AttackWorker
	lock            sync.Mutex
	errNowExecuting = errors.New("attacker is executing now")
)

// AttackWorker aggregates the instances to attack
type AttackWorker struct {
	attacker *vegeta.Attacker
	targeter vegeta.Targeter
	rate     uint64
	duration time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// StopAttack stop current attack
func StopAttack() bool {
	lock.Lock()
	defer lock.Unlock()

	if currentWorker == nil {
		return false
	}
	// cancel jobs
	currentWorker.cancel()

	return true
}

// NewAttackWorker returns a new AttackWorker with default options
func NewAttackWorker(atk *vegeta.Attacker, tr vegeta.Targeter, rate uint64, duration time.Duration) *AttackWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &AttackWorker{
		attacker: atk,
		targeter: tr,
		rate:     rate,
		duration: duration,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Run register vegeta attack job
func (worker *AttackWorker) Run(resultsBasePath string) error {
	lock.Lock()
	defer lock.Unlock()

	if currentWorker != nil {
		return errNowExecuting
	}

	// do attack
	currentWorker = worker
	go func() {
		if err := worker.attack(resultsBasePath, webSocketHub); err != nil {
			log.Println(err)
		}
		worker.UnbindWorker()
	}()

	return nil
}

// UnbindWorker remove current worker binding
func (worker *AttackWorker) UnbindWorker() {
	lock.Lock()
	defer lock.Unlock()
	if currentWorker == worker {
		currentWorker = nil
	}
}

func (worker *AttackWorker) attack(baseFilePath string, broadcaster Broadcaster) error {
	tmpFilePath := tmpFile(baseFilePath)

	out, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %s", tmpFilePath, err)
	}
	defer out.Close()

	atk := worker.attacker
	res := atk.Attack(worker.targeter, worker.rate, worker.duration)
	enc := vegeta.NewEncoder(out)

	log.Println("start attack", tmpFilePath)
	broadcaster.Broadcast("attackStart", map[string]interface{}{
		"rate":     worker.rate,
		"duration": worker.duration,
	})

	metrics := &vegeta.Metrics{}
	defer metrics.Close()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

attack:
	for {
		select {
		case <-worker.ctx.Done():
			atk.Stop()
			log.Println("stopped attack", tmpFilePath)
			return nil
		case <-ticker.C:
			broadcaster.Broadcast("attackMetrics", metrics)
		case r, ok := <-res:
			if !ok {
				broadcaster.Broadcast("attackMetrics", metrics)
				break attack
			}
			if err = enc.Encode(r); err != nil {
				return err
			}
			// add result to report
			metrics.Add(r)
			metrics.Close()
			log.Println(r)
		}
	}

	// from 00000.tmp to 00000.bin
	out.Close()
	resultFilePath := resultFile(baseFilePath)
	if err := os.Rename(tmpFilePath, resultFilePath); err != nil {
		return err
	}

	paths := strings.Split(resultFilePath, "/")
	broadcaster.Broadcast("attackFinish", map[string]interface{}{
		"filename": paths[len(paths) - 1],
	})

	log.Println("finish attack", tmpFilePath)

	return nil
}

func tmpFile(basePath string) string {
	filename := fmt.Sprintf("%d.tmp", time.Now().Unix())
	return filepath.Join(basePath, filename)
}

func resultFile(basePath string) string {
	filename := fmt.Sprintf("%d.bin", time.Now().Unix())
	return filepath.Join(basePath, filename)
}
