package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{
		name:               "home",
		url:                "/",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "about",
		url:                "/about",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "generals-quarters",
		url:                "/generals-quarters",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "majors-suites",
		url:                "/majors-suites",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "search-availability",
		url:                "/search-availability",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "contact",
		url:                "/contact",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "make-reservation",
		url:                "/make-reservation",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	}, {
		name:   "post-search-availability",
		url:    "/search-availability",
		method: "POST",
		params: []postData{
			{
				key:   "start",
				value: "2020-01-01",
			}, {
				key:   "end",
				value: "2020-01-02",
			},
		},
		expectedStatusCode: http.StatusOK,
	}, {
		name:   "post-search-availability-json",
		url:    "/search-availability-json",
		method: "POST",
		params: []postData{
			{
				key:   "start",
				value: "2020-01-01",
			}, {
				key:   "end",
				value: "2020-01-02",
			},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "make-reservation-post",
		url:    "/make-reservation",
		method: "POST",
		params: []postData{
			{
				key:   "firstName",
				value: "Majed",
			}, {
				key:   "lastName",
				value: "Nabavian",
			}, {
				key:   "email",
				value: "majedutd@gmail.com",
			}, {
				key:   "phone",
				value: "555-555-5555",
			},
		},
		expectedStatusCode: http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//	create a test server
	ts := httptest.NewTLSServer(routes)
	//	we close this server when we are done with it
	// defer close the ts after the testHandlers is finished
	// defer actually makes the closing parts wait
	defer ts.Close()
	for _, test := range theTests {
		if test.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal()
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d status code!", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal()
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d status code!", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
