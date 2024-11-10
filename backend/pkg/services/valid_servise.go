package services

import (
	"project-root/pkg/api"
	"regexp"
	"strings"
	"time"
)

type ValidatorFunc func(interface{}) []api.ValidationError

// Validators is a map of operation types to their respective validation functions
var Validators = map[string]ValidatorFunc{
	"registration": validateRegistration,
	"post":         validatePost,
	"comment":      validateComment,
}

// ValidateOperation validates the operation based on operationType and data.
func ValidateOperation(operationType string, data interface{}) []api.ValidationError {
	validator, ok := Validators[operationType]
	if !ok {
		return []api.ValidationError{{Field: "", Message: "Unsupported operation type"}}
	}

	// Call the validator function with type assertion
	return validator(data)
}

func containsWhitespace(s string) bool {
	return strings.Contains(s, " ")
}

func validateRegistration(data interface{}) []api.ValidationError {
	registrForm, ok := data.(api.RegistrationRequest)
	if !ok {
		return []api.ValidationError{{Field: "", Message: "Invalid data type for registration"}}
	}

	var validationErrors []api.ValidationError

	if len(registrForm.Nickname) < 3 {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "nickname",
			Message: "Nickname must be at least 3 characters long",
		})
	}
	if containsWhitespace(registrForm.Nickname) {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "nickname",
			Message: "Nickname must not contain whitespace",
		})
	}
	if !isValidEmail(registrForm.Email) {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}
	if !isValidPassword(registrForm.Password) {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "password",
			Message: "Password must be at least 6 characters long and contain at least one uppercase letter, one lowercase letter, one digit, and one special character",
		})
	}
	if registrForm.Age == "" {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "age",
			Message: "Age is required",
		})
	} else {
		if !isValidDate(registrForm.Age) {
			validationErrors = append(validationErrors, api.ValidationError{
				Field:   "age",
				Message: "Invalid date format. Date must be in dd.mm.yyyy format",
			})
		} else {
			ageDate, _ := time.Parse("02.01.2006", registrForm.Age)
			today := time.Now()
			if ageDate.After(today) {
				validationErrors = append(validationErrors, api.ValidationError{
					Field:   "age",
					Message: "Date of birth cannot be in the future",
				})
			}
		}
	}
	if registrForm.FirstName == "" {
		message := "First name is required"
		if containsWhitespace(registrForm.FirstName) {
			message += " and must not contain whitespace"
		}
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "first_name",
			Message: message,
		})
	}
	if registrForm.LastName == "" {
		message := "Last name is required"
		if containsWhitespace(registrForm.LastName) {
			message += " and must not contain whitespace"
		}
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "last_name",
			Message: message,
		})
	}
	if registrForm.Gender == "" {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "gender",
			Message: "Gender is required",
		})
	}
	return validationErrors
}

// isValidPassword checks if the provided password meets the required criteria.
func isValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsAny(string(char), "!@#$%^&*()-_=+{}[]|/<>?,."):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// isValidEmail checks if the provided email is in a valid format.
func isValidEmail(email string) bool {
	// Regular expression for basic email validation
	// This is a simple example and may not cover all valid email formats
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidDate checks if the provided date string is in dd.mm.yyyy format and represents a valid date.
func isValidDate(dateStr string) bool {
	_, err := time.Parse("02.01.2006", dateStr)
	return err == nil
}

// validatePost validates the fields of PostCreateRequest.
// Returns a slice of ValidationError if any field is invalid.
func validatePost(data interface{}) []api.ValidationError {
	postData, ok := data.(api.PostCreateRequest)
	if !ok {
		return []api.ValidationError{{Field: "", Message: "Invalid data type for post"}}
	}

	var validationErrors []api.ValidationError

	// Validate title length (6 to 48 characters)
	if len(postData.Title) < 6 || len(postData.Title) > 48 {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "title",
			Message: "Title length must be between 6 and 48 characters",
		})
	}

	// Trim leading and trailing whitespace from content
	postData.Content = strings.TrimSpace(postData.Content)

	// Validate content length (0 to 256 characters)
	if len(postData.Content) > 256 {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "content",
			Message: "Content length cannot exceed 256 characters",
		})
	}

	// Validate categories
	if postData.Categories == "" {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "categories",
			Message: "Categories are required",
		})
	} else {
		categories := ParseCategories(postData.Categories)
		if len(categories) < 1 || len(categories) > 5 {
			validationErrors = append(validationErrors, api.ValidationError{
				Field:   "categories",
				Message: "There must be between 1 and 5 categories",
			})
		}
		for _, category := range categories {
			if len(category) < 1 || len(category) > 15 {
				validationErrors = append(validationErrors, api.ValidationError{
					Field:   "categories",
					Message: "Each category must be between 1 and 10 characters long",
				})
			}
			if containsWhitespace(category) {
				validationErrors = append(validationErrors, api.ValidationError{
					Field:   "categories",
					Message: "Categories cannot contain whitespace between words",
				})
			}
		}
	}

	return validationErrors
}

// validateComment validates the fields of CommentCreateRequest.
// Returns a slice of ValidationError if any field is invalid.
func validateComment(data interface{}) []api.ValidationError {
	commentData, ok := data.(api.CommentCreateRequest)
	if !ok {
		return []api.ValidationError{{Field: "", Message: "Invalid data type for comment"}}
	}

	var validationErrors []api.ValidationError

	// Trim leading and trailing whitespace from content
	commentData.Content = strings.TrimSpace(commentData.Content)

	// Validate content length (1 to 256 characters)
	if len(commentData.Content) < 1 || len(commentData.Content) > 256 {
		validationErrors = append(validationErrors, api.ValidationError{
			Field:   "content",
			Message: "Content must be between 1 and 256 characters long",
		})
	}

	return validationErrors
}
