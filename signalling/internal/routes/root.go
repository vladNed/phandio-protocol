package routes

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mvx-mnr-atomic/signalling/internal/hub"
	"github.com/mvx-mnr-atomic/signalling/internal/logging"
	"github.com/mvx-mnr-atomic/signalling/internal/router"
)

var (
	upgrader   = websocket.Upgrader{}
	logger     = logging.GetLogger(nil)
	HttpRouter = router.NewHttpRouter("/api/v1", []router.Route{
		{
			Path:        "/ping",
			Method:      http.MethodGet,
			HandlerFunc: ping,
		},
	})
	WSRouter = router.NewWsRouter("/ws/v1", []router.Route{
		{
			Path: "/",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				serveWebSocket(hub.HubInstance, w, r)
			},
		},
	})
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
