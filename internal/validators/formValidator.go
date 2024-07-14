package validators

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type FormValidator struct {
	GenericErrors []string
	FieldErrors   map[string]string
}

func (v *FormValidator) Valid() bool {
	return len(v.GenericErrors) == 0 && len(v.FieldErrors) == 0
}

func (v *FormValidator) AddGenericError(message string) {
	v.GenericErrors = append(v.GenericErrors, message)
}

func (v *FormValidator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *FormValidator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func InBetween(value, low, high int) bool {
	return low <= value && value <= high
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
