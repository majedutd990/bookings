package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	myH := myHandler{}
	h := NoSurf(&myH)
	//	we should test that what it returns actually ia a http handler

	switch v := h.(type) {
	case http.Handler:
	//	do nothing
	default:
		t.Error(fmt.Sprintf("type is not http handler noSurf(). but is %T\n.", v))
	}
}

func TestSessionLoad(t *testing.T) {
	myH := myHandler{}
	h := SessionLoad(&myH)
	//	we should test that what it returns actually ia a http handler

	switch v := h.(type) {
	case http.Handler:
	//	do nothing
	default:
		t.Error(fmt.Sprintf("type is not http handler sessionLoad(). but is %T\n.", v))
	}
}
