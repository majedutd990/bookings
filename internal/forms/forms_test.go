package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestForm_Valid is like this because it has a receiver
func TestForm_Valid(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid where should have found valid")
	}
}

func TestForm_Required(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form should have been invalid but here is valid Required()")
	}
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	form = New(postData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form should have been valid but here is invalid Required()")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postValues := url.Values{}
	form := New(postValues)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("IsEmail. form is valid for a non existing field")
	}
	// we also check Get function of error file here too
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}
	postValues = url.Values{}
	postValues.Add("emailTest", "test@gmail.com")
	form = New(postValues)
	form.IsEmail("emailTest")
	if !form.Valid() {
		t.Error("IsEmail. input is email but return false")
	}
	postValues = url.Values{}
	postValues.Add("emailTest2", "test")
	form = New(postValues)
	form.IsEmail("emailTest2")
	if form.Valid() {
		t.Error("IsEmail. input is not email return true")
	}
}

func TestForm_MinLength(t *testing.T) {
	postValues := url.Values{}
	form := New(postValues)
	form.MinLength("emptyTest", 10)
	if form.Valid() {
		t.Error("MinLength works for an empty field")
	}
	postValues = url.Values{}
	postValues.Add("test", "tests")
	form = New(postValues)
	form.MinLength("test", 4)
	if !form.Valid() {
		t.Error("minlength. shows false although requirement is met")
	}
	// we also check Get function of error file here too
	isError := form.Errors.Get("test")
	if isError != "" {
		t.Error("should not have an error but get one")
	}
	postValues = url.Values{}
	postValues.Add("test2", "tests")
	form = New(postValues)
	form.MinLength("test2", 500)
	if form.Valid() {
		t.Error("minlength. show min length of 500 when data is shorter")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	has := form.Has("emptyTest")
	if has {
		t.Error("Has. form does not have the field but return true")
	}

	postData = url.Values{}
	postData.Add("test", "test")
	form = New(postData)
	has = form.Has("test")
	if !has {
		t.Error("Has. form has field but return false")
	}
}
