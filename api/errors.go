package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

// Error implements the Error interface
func (e Error) Error() string {
	return e.Err
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id give",
	}
}

func ErrUnAuthorizedD() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "inalid JSON request",
	}
}

func ErrResoruceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + "resource not found",
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}
