package domain

import (
	"github.com/gofiber/fiber/v2"
)

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// HTTPErrorValidation example
type HTTPErrorValidation struct {
	Field    string    `json:"field" example:"username"`
	Message string `json:"message" example:"username cannot empty"`
}

// NewHttpError example
func NewHttpError(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(HTTPError{Code: status, Message: err.Error()})
}
