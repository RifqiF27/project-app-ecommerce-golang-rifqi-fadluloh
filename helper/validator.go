package helper

import (
	"ecommerce/model"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("regex", func(fl validator.FieldLevel) bool {
		pattern := fl.Param()
		value := fl.Field().String()

		matched, err := regexp.MatchString(pattern, value)
		if err != nil {
			return false
		}
		return matched
	})

	return &Validator{validate: v}
}

func (v *Validator) ValidateStruct(data interface{}) error {
	return v.validate.Struct(data)
}
func (v *Validator) ValidateLoginStruct(req model.User) error {
	return v.validate.StructExcept(req, "Name")
}

func FormatValidationError(err error) string {
	errorMessages := map[string]string{
		"Name_required":          "Name is required",
		"Name_min":               "Name must have at least 3 characters",
		"Name_regex":             "Name must only contain letters and spaces",
		"Email_required_without": "Email or Phone is required",
		"Email_email":            "Email format is invalid",
		"Phone_required_without": "Phone or Email is required",
		"Phone_min":              "Phone number must be at least 10 digits",
		"Phone_max":              "Phone number must be at most 13 digits",
		"Phone_numeric":          "Phone number must be numeric",
		"Password_required":      "Password is required",
		"Password_min":           "Password must have at least 8 characters",
	}

	var errMessages []string
	for _, err := range err.(validator.ValidationErrors) {

		key := err.Field() + "_" + err.Tag()

		if message, found := errorMessages[key]; found {
			errMessages = append(errMessages, message)
		} else {
			errMessages = append(errMessages, err.Field()+" is invalid: "+err.Tag())
		}
	}

	return strings.Join(errMessages, ", ")

}
