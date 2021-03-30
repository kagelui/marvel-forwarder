package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kagelui/marvel-forwarder/cmd/serverd/handler"
	"github.com/kagelui/marvel-forwarder/internal/pkg/server"
	"github.com/kagelui/marvel-forwarder/internal/service/characters"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("server started")

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err.Error())
		os.Exit(132)
	}

	modelStore := &characters.ModelStore{DB: db}

	mux := http.NewServeMux()
	mux.Handle("/characters", handler.WrapError(handler.GetMarvelCharacterList(modelStore)))
	mux.Handle("/characters/{id}", handler.WrapError(handler.GetMarvelCharacterDetail))

	server.New(":8080", mux).Start()
}
