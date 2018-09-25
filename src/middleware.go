package main

import (
	"errors"
	"net/http"

	"github.com/tomasen/realip"
)

func universalLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logger.UniversalLogger(realip.FromRequest(r), r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func isAuth(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	h := func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("auth")
		if c.UserConfig.RequireAuth {
			if auth != c.UserConfig.AuthKey && !c.UserConfig.Debug {
				errorResponse(w, errors.New("unauthorized"))
				return
			}
		}
		handler(w, r)
	}
	return h
}

func isConfig(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if !c.Config && !c.UserConfig.Debug {
			errorResponse(w, errors.New("service not configured"))
			return
		}
		handler(w, r)
	}
	return h
}
