package main

import (
  "net/http"

  "github.com/labstack/echo"
	vegeta "github.com/tsenart/vegeta/lib"
)

// Handler is the HTTP handler
type Handler struct {}

// Index handle GET / request
func (h *Handler) Index(c echo.Context) error {
  return c.String(http.StatusOK, "Hello Trunks!!")
}

// PostAttack handle POST /attack request
func (h *Handler) PostAttack(c echo.Context) error {
	opts := &AttackOptions{
		Headers: headers{http.Header{}},
		Laddr: localAddr{&vegeta.DefaultLocalAddr},
	}

	if err := c.Bind(opts); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

  return c.JSON(http.StatusOK, opts)
}
