package api

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body,omitempty"`
}
type ValidationError struct {
	FieldName    string `json:"fieldName"`
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
}

func SuccessResponse(body interface{}) Response {
	return Response{
		Code:    "000",
		Message: "Success",
		Body:    body,
	}
}

func Err(code string, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
func ErrWithBody(code string, message string, body interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Body:    body,
	}
}
func JwtError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Err("401", message))

}

func Unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Err("401", fiber.ErrUnauthorized.Error()))
}

func Forbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(Err("403", fiber.ErrForbidden.Error()))
}

func InternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Err("500", fmt.Sprintf("%s", message)))
}

func Ok(c *fiber.Ctx, body interface{}) error {
	return c.Status(fiber.StatusOK).JSON(SuccessResponse(body))
}

func NotFound(c *fiber.Ctx, message string) error {
	msg := fiber.ErrBadRequest.Message
	if message != "" {
		msg = message
	}
	return c.Status(fiber.StatusNotFound).JSON(Err("404", msg))
}

func BadRequest(c *fiber.Ctx, message string) error {
	msg := fiber.ErrBadRequest.Message
	if message != "" {
		msg = message
	}
	return c.Status(fiber.StatusBadRequest).JSON(Err("400", msg))
}
func ValidateError(c *fiber.Ctx, body interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrWithBody("400", "Invalid Request", body))
}

func ValidationErrorResponse(c *fiber.Ctx, err error, request interface{}) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := []ValidationError{}

		for _, fe := range ve {
			jsonField := GetJSONFieldName(request, fe.StructField())
			switch fe.Tag() {
			case "required":
				errs = append(errs, ValidationError{
					FieldName:    jsonField,
					ErrorType:    "required",
					ErrorMessage: fe.Error(),
				})
			case "min":
				errs = append(errs, ValidationError{
					FieldName:    jsonField,
					ErrorType:    "min",
					ErrorMessage: fe.Error(),
				})
			case "max":
				errs = append(errs, ValidationError{
					FieldName:    jsonField,
					ErrorType:    "max",
					ErrorMessage: fe.Error(),
				})
			default:
				errs = append(errs, ValidationError{
					FieldName:    jsonField,
					ErrorType:    fe.Tag(),
					ErrorMessage: fe.Error(),
				})
			}
		}
		return ValidateError(c, errs)
	}

	return BadRequest(c, "Invalid request")
}

func GetJSONFieldName(val interface{}, fieldName string) string {
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if f, ok := t.FieldByName(fieldName); ok {
		tag := f.Tag.Get("json")
		if tag != "" && tag != "-" {
			if idx := indexComma(tag); idx >= 0 {
				return tag[:idx]
			}
			return tag
		}
	}
	return fieldName
}

func indexComma(s string) int {
	for i, c := range s {
		if c == ',' {
			return i
		}
	}
	return -1
}
