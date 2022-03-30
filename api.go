package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type api struct {
	store Store
}

// ctxKey used for create a request context to getParam from path.
type ctxKey struct{}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler func(h *api, w http.ResponseWriter, r *http.Request)
}

// Serve responsible for fulfill requests.
func (h *api) Serve(w http.ResponseWriter, r *http.Request) {
	routes := []route{
		{
			method:  "GET",
			regex:   regexp.MustCompile("^/ping$"),
			handler: pong,
		},
		{
			method:  "GET",
			regex:   regexp.MustCompile("^/([^/]+)$"),
			handler: get,
		},
		{
			method:  "POST",
			regex:   regexp.MustCompile("^/([^/]+)$"),
			handler: set,
		},
		{
			method:  "PUT",
			regex:   regexp.MustCompile("^/([^/]+)$"),
			handler: set,
		},
		{
			method:  "PATCH",
			regex:   regexp.MustCompile("^/([^/]+)$"),
			handler: set,
		},
		{
			method:  "DELETE",
			regex:   regexp.MustCompile("^/flush$"),
			handler: flush,
		},
		{
			method:  "DELETE",
			regex:   regexp.MustCompile("^/([^/]+)$"),
			handler: drop,
		},
	}

	log.Printf("[http] %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	var methods []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				methods = append(methods, route.method)

				continue
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])

			w.Header().Set("Content-Type", "application/json")
			route.handler(h, w, r.WithContext(ctx))

			return
		}
	}

	if len(methods) > 0 {
		w.Header().Set("Allow", strings.Join(methods, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func getParam(r *http.Request) string {
	fields, _ := r.Context().Value(ctxKey{}).([]string)

	return fields[0]
}

func pong(_ *api, w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "{\"key\":\"%s\",\"value\":\"%s\"}", "ping", "fooooooooooo")
}

func set(h *api, w http.ResponseWriter, r *http.Request) {
	key := getParam(r)
	val := r.FormValue("value")

	if val == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	h.store.Set(key, val)

	_, _ = fmt.Fprintf(w, "{\"key\":\"%s\",\"value\":\"%s\"}", key, val)
}

func get(h *api, w http.ResponseWriter, r *http.Request) {
	key := getParam(r)

	val, found := h.store.Get(key)

	if found {
		_, _ = fmt.Fprintf(w, "{\"key\":\"%s\",\"value\":\"%s\"}", key, val)

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func drop(h *api, w http.ResponseWriter, r *http.Request) {
	key := getParam(r)

	h.store.Delete(key)

	w.WriteHeader(http.StatusAccepted)
}

func flush(h *api, w http.ResponseWriter, _ *http.Request) {
	h.store.Flush()

	w.WriteHeader(http.StatusAccepted)
}
