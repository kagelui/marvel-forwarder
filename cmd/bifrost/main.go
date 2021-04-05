package main

import (
	"context"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kagelui/marvel-forwarder/internal/pkg/envvar"
	"github.com/kagelui/marvel-forwarder/internal/pkg/loglib"
	"github.com/kagelui/marvel-forwarder/internal/service/marvel"
	_ "github.com/lib/pq"
)

const retries = 3

func main() {
	lg := loglib.DefaultLogger()
	ctx := loglib.SetLogger(context.Background(), lg)

	lg.InfoF("starting syncing with marvel API...")

	var e envVar

	if err := envvar.Read(&e); err != nil {
		lg.ErrorF(err.Error())
		os.Exit(1)
	}

	client := marvel.ApiClient{
		Client:     http.DefaultClient,
		PublicKey:  e.PublicKey,
		PrivateKey: e.PrivateKey,
		APIAddr:    e.APIAddr,
		Retries:    retries,
	}

	db, err := sqlx.Connect("postgres", e.DBAddr)
	if err != nil {
		lg.ErrorF(err.Error())
		os.Exit(2)
	}

	characters, err := client.RetrieveCharacters(ctx)
	if err != nil {
		lg.ErrorF(err.Error())
		os.Exit(3)
	}

	if err = characters.Save(ctx, db); err != nil {
		lg.ErrorF(err.Error())
		os.Exit(4)
	}
}

type envVar struct {
	PublicKey  string `env:"PUBLIC_KEY"`
	PrivateKey string `env:"PRIVATE_KEY"`
	APIAddr    string `env:"MARVEL_API_URL"`
	DBAddr     string `env:"DATABASE_URL"`
}
