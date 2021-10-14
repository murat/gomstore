package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_api_Serve(t *testing.T) {
	type req struct {
		method, url, body string
	}
	type resp struct {
		status uint
		body   string
	}
	requests := []struct {
		title    string
		request  req
		response *resp
	}{
		{
			title: "not found route",
			request: req{
				method: "HEAD",
				url:    "/",
			},
			response: &resp{
				status: http.StatusNotFound,
				body:   "",
			},
		},
		{
			title: "happy ping",
			request: req{
				method: "GET",
				url:    "/ping",
			},
			response: &resp{
				status: http.StatusOK,
				body:   `{"key":"ping","value":"pong"}`,
			},
		},
		{
			title: "bad ping",
			request: req{
				method: "HEAD",
				url:    "/ping",
			},
			response: &resp{
				status: http.StatusMethodNotAllowed,
				body:   "",
			},
		},
		{
			title: "happy set",
			request: req{
				method: "POST",
				url:    "/foo",
				body:   "value=bar",
			},
			response: &resp{
				status: http.StatusOK,
				body:   `{"key":"foo","value":"bar"}`,
			},
		},
		{
			title: "bad set",
			request: req{
				method: "POST",
				url:    "/foo",
				body:   "",
			},
			response: &resp{
				status: http.StatusUnprocessableEntity,
				body:   "",
			},
		},
		{
			title: "happy get",
			request: req{
				method: "GET",
				url:    "/foo",
			},
			response: &resp{
				status: http.StatusOK,
				body:   `{"key":"foo","value":"bar"}`,
			},
		},
		{
			title: "bad get",
			request: req{
				method: "GET",
				url:    "/not-exist-key",
			},
			response: &resp{
				status: http.StatusNotFound,
				body:   "",
			},
		},
		{
			title: "happy delete",
			request: req{
				method: "DELETE",
				url:    "/foo",
			},
			response: &resp{
				status: http.StatusAccepted,
				body:   "",
			},
		},
		{
			title: "happy flush",
			request: req{
				method: "DELETE",
				url:    "/flush",
			},
			response: &resp{
				status: http.StatusAccepted,
				body:   "",
			},
		},
	}

	for _, tt := range requests {
		api := &api{
			store: NewStore(),
		}
		api.store.Set("foo", "bar")

		t.Run(tt.title, func(t *testing.T) {
			var data io.Reader
			if tt.request.body != "" {
				data = strings.NewReader(tt.request.body)
			}

			request, err := http.NewRequest(tt.request.method, tt.request.url, data)
			if err != nil {
				t.Fatal(err)
			}

			if tt.request.body != "" {
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(api.Serve)

			handler.ServeHTTP(recorder, request)

			if status := recorder.Code; status != int(tt.response.status) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.response.status)
			}

			expected := tt.response.body
			if recorder.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					recorder.Body.String(), expected)
			}
		})
	}
}
