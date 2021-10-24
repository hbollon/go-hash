package gohash

import "github.com/go-playground/validator/v10"

var structValidator *validator.Validate

func init() {
	structValidator = validator.New()
	structValidator.RegisterValidation("is-hash-method", ValidateHashMethod)
}

func ValidateHashMethod(fl validator.FieldLevel) bool {
	return fl.Field().String() == "MD5" || fl.Field().String() == "SHA1"
}
