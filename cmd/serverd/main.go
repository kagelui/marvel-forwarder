package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.Handle("/characters", handler.WrapError(handler.GetMarvelCharacterList(modelStore))).Methods("GET")
	r.Handle("/characters/{id:[0-9]+}", handler.WrapError(handler.GetMarvelCharacterDetail(modelStore))).Methods("GET")

	server.New(":8080", r).Start()
}
