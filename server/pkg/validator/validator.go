package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	// Register JSON tag name function to use JSON tag names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(structure interface{}) []ValidationError {
	err := validate.Struct(structure)
	if err == nil {
		return nil
	}

	var errors []ValidationError

	for _, err := range err.(validator.ValidationErrors) {
		var element ValidationError
		// Use the field name from the tag name function (which returns JSON tag name)
		element.Field = err.Field()
		element.Message = msgForTag(err.Tag(), err.Param())
		errors = append(errors, element)
	}

	return errors
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is: %s", param)
	case "max":
		return fmt.Sprintf("Maximum length is: %s", param)
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", strings.Join(strings.Split(param, " "), ", "))
	case "url":
		return "Invalid url format"
	case "gte":
		return fmt.Sprintf("Must be greater than or equal: %s", param)
	default:
		return fmt.Sprintf("Failed on tag: %s", tag)
	}
}
