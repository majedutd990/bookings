package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

//Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//New Initializes a Form struct
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: errors(map[string][]string{}),
	}
}

//Required has a variadic function that take as many arg as it can and checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank!")
		}
	}
}

//Has checks if field is in Form and not empty
func (f *Form) Has(field string, r *http.Request) bool {

	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

//MinLength checks for minimum length of string fields
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

//IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
