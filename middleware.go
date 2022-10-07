package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func parsedBody() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do stuff
			var currentBody car
			body, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(body, &currentBody)
			r.WithContext(context.WithValue(r.Context(), "jbody", currentBody))
			h.ServeHTTP(w, r)
		})
	}
}
