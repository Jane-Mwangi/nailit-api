package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Define a new Validator type which contains a map of validation errors.
type Validator struct {
	Errors map[string]string
}

// New is a helper which creates a new validator instance.
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// valid returns a true if the validator has no errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adss an error meassge to themap (so long as no entry already exists for the given key)
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// check  adds an error message to the map only if the  validator  check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if a specific value is in a list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// matches returns true is all the string value matches a specifc regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

//unique returns true if al the string values in a slice are unique

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
