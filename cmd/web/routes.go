package web

import (
	"WebsocketPractice/internal/chat"
	"github.com/bmizerany/pat"
	"net/http"
)

func Routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(chat.Home))
	mux.Get("/ws", http.HandlerFunc(chat.WsEndpoint))
	return mux
}
