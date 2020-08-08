package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// CustomValidator contains validator
type CustomValidator struct {
	validator *validator.Validate
}

// New creates and returns a new validator
func New() *CustomValidator {
	ctm := &CustomValidator{validator: validator.New()}
	ctm.registerCustomValidation()

	return ctm
}

// Validate validates given struct
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	return customError(err)
}

func customError(err error) error {
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				return NewValidationError(fmt.Sprintf("%s is required",
					err.Field()))
			case "email":
				return NewValidationError(fmt.Sprintf("%s is not valid email",
					err.Field()))
			case "gte":
				return NewValidationError(fmt.Sprintf("%s value must be greater than %s",
					err.Field(), err.Param()))
			case "lte":
				return NewValidationError(fmt.Sprintf("%s value must be lower than %s",
					err.Field(), err.Param()))
			case "unique":
				return NewValidationError(fmt.Sprintf("%s is already taken",
					err.Field()))
			default:
				return NewValidationError(fmt.Sprintf("%s validation error on %s tag", err.Field(), err.ActualTag()))
			}
		}
	}
	return nil
}

func (cv *CustomValidator) registerCustomValidation() {
	_ = cv.validator.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		// TODO: validate unique field
		return true
	})
}

type ValidationErrors struct {
	err error
}

func (v ValidationErrors) Error() string {
	return v.err.Error()
}

func NewValidationError(msg string) ValidationErrors {
	return ValidationErrors{errors.New(msg)}
}
