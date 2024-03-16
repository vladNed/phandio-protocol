package router

import (
	"fmt"
	"net/http"

	"github.com/mvx-mnr-atomic/signalling/internal/logging"
)

var logger = logging.GetLogger(nil)

type HttpRouter struct {
	prefix	string
	routes 	[]Route
}

func NewHttpRouter(prefix string, routes []Route) *HttpRouter {
	return &HttpRouter{
		prefix: prefix,
		routes: routes,
	}
}

func (r *HttpRouter) Register() {
	for _, route := range r.routes {
		routePath := r.prefix + route.Path
		logger.Info(fmt.Sprintf("HTTP Route -> %s [%s]", routePath, route.Method))
		http.HandleFunc(routePath, HttpMethodHandler(route.Method, route.HandlerFunc))
	}
}

type WsRouter struct {
	prefix	string
	routes 	[]Route
}

func NewWsRouter(prefix string, routes []Route) *WsRouter {
	return &WsRouter{
		prefix: prefix,
		routes: routes,
	}
}

func (r *WsRouter) Register() {
	for _, route := range r.routes {
		routePath := r.prefix + route.Path
		logger.Info(fmt.Sprintf("WS Route -> %s", routePath))
		http.HandleFunc(routePath, route.HandlerFunc)
	}
}