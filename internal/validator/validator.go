package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Use the regexp.MustCompile() function to parse a regular expression pattern
// for sanity checking the format of an email address. This returns a pointer to
// a 'compiled' regexp.Regexp type, or panics in the event of an error. Parsing
// this pattern once at startup and storing the compiled *regexp.Regexp in a
// variable is more performant than re-parsing the pattern each time we need it.
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z]{2,})+$`)

// defie a validator type that contains any fielderrors
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid returns true if the fielderrors map don't contain any errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
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

// Adds a NonfieldError to the string slice
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

// returns true if char counts is equal or more than n
func MinCharCount(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
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

// return true if a string matches a provided compiled regular expressio pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
