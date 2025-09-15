package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Enum interface {
	IsValid() bool
}

func validateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)

	return value.IsValid()
}

func RegisterValidators() {
	validate = validator.New()

	err := validate.RegisterValidation("enum", validateEnum)

	if err != nil {
		panic(err)
	}
}

func Struct(s any) error {
	return validate.Struct(s)
}
