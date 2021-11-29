package main

import (
	"WebsocketPractice/cmd/web"
	"log"
	"net/http"
)

func main() {
	routes := web.Routes()
	log.Println("Starting server on por 8080")
	_ = http.ListenAndServe(":8080", routes)

}
