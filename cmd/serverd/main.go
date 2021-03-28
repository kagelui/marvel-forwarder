package main

import (
	"github.com/kagelui/marvel-forwarder/cmd/serverd/handler"
	"log"
	"net/http"

	"github.com/kagelui/marvel-forwarder/internal/pkg/server"
)

func main() {
	log.Println("server started")

	mux := http.NewServeMux()
	mux.Handle("/characters", handler.WrapError(handler.GetMarvelCharacterList))
	mux.Handle("/characters/{id}", handler.WrapError(handler.GetMarvelCharacterDetail))

	server.New(":8080", mux).Start()
}
