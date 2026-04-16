package httpx

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Status  int               `json:"-"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func NewFieldError(status int, message string, fields map[string]string) *AppError {
	return &AppError{Status: status, Message: message, Fields: fields}
}

func BadRequest(message string) *AppError {
	return NewError(fiber.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return NewError(fiber.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return NewError(fiber.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return NewError(fiber.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return NewError(fiber.StatusConflict, message)
}

func Internal(message string) *AppError {
	return NewError(fiber.StatusInternalServerError, message)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		payload := fiber.Map{
			"success":    false,
			"message":    appErr.Message,
			"request_id": requestID(c),
		}
		if len(appErr.Fields) > 0 {
			payload["fields"] = appErr.Fields
		}
		return c.Status(appErr.Status).JSON(payload)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"success":    false,
			"message":    fiberErr.Message,
			"request_id": requestID(c),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success":    false,
		"message":    "internal server error",
		"request_id": requestID(c),
	})
}
