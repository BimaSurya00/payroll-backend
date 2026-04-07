package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"hris/shared/helper"
)

var validate *validator.Validate

func InitValidator() *validator.Validate {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("password_strength", passwordStrength)
	_ = validate.RegisterValidation("trimmed_string", trimmedString)
	_ = validate.RegisterValidation("time_format", timeFormat)
	_ = validate.RegisterValidation("slug", slugValidator)

	return validate
}

func GetValidator() *validator.Validate {
	if validate == nil {
		return InitValidator()
	}
	return validate
}

func ValidateStruct(s any) []helper.ValidationError {
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
	case "time_format":
		return "Invalid time format. Use HH:MM format (e.g., 09:00, 17:30)"
	case "slug":
		return "Slug can only contain lowercase letters, numbers, and hyphens"
	case "regexp":
		if err.Param() == "^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$$" {
			return "Invalid time format. Use HH:MM format (e.g., 09:00, 17:30)"
		}
		return fmt.Sprintf("%s is invalid", err.Field())
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

func timeFormat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Validate HH:MM format (24-hour)
	matched, _ := regexp.MatchString(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`, value)
	return matched
}

func slugValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Validate slug: lowercase letters, numbers, and hyphens only
	matched, _ := regexp.MatchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, value)
	return matched
}
