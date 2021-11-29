package main

import (
	"WebsocketPractice/cmd/web"
	"WebsocketPractice/internal/chat"
	"fmt"
	"log"
	"net/http"
)

func main() {
	routes := web.Routes()

	go chat.ListenToWsChannel()
	fmt.Println("Started listening to messages")

	log.Println("Starting server on port 8080")
	_ = http.ListenAndServe(":8080", routes)
}
