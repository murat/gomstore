package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type api struct {
	store *store
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler func(h *api, w http.ResponseWriter, r *http.Request)
}

// Serve responsible for fulfill requests
func (h *api) Serve(w http.ResponseWriter, r *http.Request) {
	routes := []route{
		{
			method:  "GET",
			regex:   regexp.MustCompile("^/ping$"),
			handler: pong,
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

			w.Header().Set("Content-Type", "application/json")
			route.handler(h, w, r)
			return
		}
	}

	if len(methods) > 0 {
		w.Header().Set("Allow", strings.Join(methods, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func pong(_ *api, w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "{\"%s\": \"%s\"}", "ping", "pong")
}
