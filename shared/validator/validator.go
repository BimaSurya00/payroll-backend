package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
)

var validate *validator.Validate

func InitValidator() *validator.Validate {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("password_strength", passwordStrength)
	_ = validate.RegisterValidation("trimmed_string", trimmedString)

	return validate
}

func GetValidator() *validator.Validate {
	if validate == nil {
		return InitValidator()
	}
	return validate
}

func ValidateStruct(s interface{}) []helper.ValidationError {
	var errors []helper.ValidationError

	err := GetValidator().Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, helper.ValidationError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
			})
		}
	}

	return errors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "password_strength":
		return "Password must contain at least one uppercase, one lowercase, one number, and one special character"
	case "trimmed_string":
		return fmt.Sprintf("%s must not have leading or trailing spaces", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// Custom validators
func passwordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func trimmedString(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == strings.TrimSpace(value) && !regexp.MustCompile(`^\s|\s$`).MatchString(value)
}
