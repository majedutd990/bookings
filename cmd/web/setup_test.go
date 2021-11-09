package main

import (
	"net/http"
	"os"
	"testing"
)

//what ever in here run before our tests
// it has to have a function called test main

func TestMain(m *testing.M) {

	// do something run the test then exit
	os.Exit(m.Run())
}

//myHandler is our own handler to do the test
type myHandler struct {
}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
