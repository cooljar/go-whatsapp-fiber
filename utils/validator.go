package utils

import (
	"github.com/cooljar/go-whatsapp-fiber/domain"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// RegisterTagNameFunc registers a function to get alternate names for StructFields.
	// eg. Title become title, CreatedAt become created_at
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) []domain.HTTPErrorValidation {
	// Define fields map.
	var fields []domain.HTTPErrorValidation

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	/*if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println(err)
		return
	}*/

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		/*fmt.Println(err.Namespace())
		fmt.Println(err.Field())
		fmt.Println(err.StructNamespace())
		fmt.Println(err.StructField())
		fmt.Println(err.Tag())
		fmt.Println(err.ActualTag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())
		fmt.Println(err.Param())*/

		fields = append(fields, domain.HTTPErrorValidation{Field: err.Field(), Message: err.Tag()})
	}

	return fields
}
