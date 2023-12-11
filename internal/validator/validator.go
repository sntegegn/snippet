package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, ok := v.FieldErrors[key]; !ok {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, count int) bool {
	return utf8.RuneCountInString(value) < count
}

func MinChars(value string, count int) bool {
	return utf8.RuneCountInString(value) > count
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func PremittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}