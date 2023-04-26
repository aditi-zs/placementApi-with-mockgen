package main

import "net/http"

func middleware(originalHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("X-API-KEY")
		contType := r.Header.Get("Content-Type")

		if val != "a601e44e306e430f8dde987f65844f05" && val != "84dcb7c09b4a4af8a67f4577ffe9b255" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("authentication failed"))

			return
		}

		if contType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = w.Write([]byte("Header Content-Type incorrect"))

			return
		}

		w.Header().Set("Content-Type", "application/json")

		originalHandler.ServeHTTP(w, r)
	}
}
