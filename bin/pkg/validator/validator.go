package validator

import (
	"fmt"
	errors "location-service/bin/pkg/http-error"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	// create custom error message like nodejs -> "message": "\"latitude\" is required"
	if cv.Validator.Struct(i) != nil {
		errs := cv.Validator.Struct(i).(validator.ValidationErrors)
		errorMsg := fmt.Sprintf("\"%s\" is %s", strings.ToLower(errs[0].Field()), errs[0].Tag())
		return errors.Conflict(errorMsg)
	}

	return nil
}

func New() *validator.Validate {
	return validator.New()
}
