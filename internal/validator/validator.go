package validator

import (
	"strings"
	"unicode/utf8"
)

// defie a validator type that contains any fielderrors
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the fielderrors map don't contain any errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFiledError function adds an error for the given key, if no error message exist for the given key
func (v *Validator) AddFiledError(key, message string) {
	// if not initialized, we will initialize the map
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Adds error message to the map when validation is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFiledError(key, message)
	}
}

// Checks if the field is blank
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// returns true if char counts is equal less than n
func MaxCharCount(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// returns true if a value is in a list of permitted int list
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
