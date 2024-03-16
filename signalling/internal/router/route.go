package router

import "net/http"

type Route struct {

	// Method is the HTTP method to be used for the route
	// NOTE: This field is empty fro websocket routes
	Method string

	// Path is the URL path to be used for the route. Must not end in trailing '/'
	Path string

	// HandlerFunc is the function to be called when the route is matched
	HandlerFunc http.HandlerFunc
}
