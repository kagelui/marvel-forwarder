package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kagelui/marvel-forwarder/cmd/serverd/handler"
	"github.com/kagelui/marvel-forwarder/internal/pkg/envvar"
	"github.com/kagelui/marvel-forwarder/internal/pkg/server"
	"github.com/kagelui/marvel-forwarder/internal/service/characters"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("server started")

	var e envVar

	if err := envvar.Read(&e); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", e.DBAddr)
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

type envVar struct {
	DBAddr string `env:"DATABASE_URL"`
}
