package main

import (
  "net/http"

  "github.com/labstack/echo"
)

// Handler is the HTTP handler
type Handler struct {}

// Index handle GET / request
func (h *Handler) Index(c echo.Context) error {
  return c.String(http.StatusOK, "Hello Trunks!!")
}
