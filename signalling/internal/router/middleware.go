package router

import "net/http"

// HttpMethodHandler is a middleware that checks if the request method is the same as the one provided
func HttpMethodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}
