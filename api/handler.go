package api

import (
	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/luganodes-task/internal/store/controller"
)

type (
	Handlers struct {
		*echo.Echo
		controller *controller.Controller
	}
)

// NewHandler creates a new handler
func NewHandler(c *controller.Controller) *Handlers {
	return &Handlers{Echo: echo.New(), controller: c}
}

// Handle handles the routes
func (h *Handlers) Handle() {
	h.GET("/deposits", h.controller.GetDeposit)
}
