package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kagelui/marvel-forwarder/internal/pkg/loglib"
	"github.com/kagelui/marvel-forwarder/internal/service/marvel"
	_ "github.com/lib/pq"
)

const retries = 3

func main() {
	lg := loglib.DefaultLogger()
	ctx := loglib.SetLogger(context.Background(), lg)

	lg.InfoF("starting syncing with marvel API...")

	e, err := readEnv()
	if err != nil {
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
	PublicKey  string
	PrivateKey string
	APIAddr    string
	DBAddr     string
}

func readEnv() (envVar, error) {
	var ev envVar

	v, ok := os.LookupEnv("MARVEL_API_URL")
	if !ok {
		return envVar{}, fmt.Errorf("missing MARVEL_API_URL")
	}
	ev.APIAddr = v
	v, ok = os.LookupEnv("PUBLIC_KEY")
	if !ok {
		return envVar{}, fmt.Errorf("missing PUBLIC_KEY")
	}
	ev.PublicKey = v
	v, ok = os.LookupEnv("PRIVATE_KEY")
	if !ok {
		return envVar{}, fmt.Errorf("missing PRIVATE_KEY")
	}
	ev.PrivateKey = v
	v, ok = os.LookupEnv("DATABASE_URL")
	if !ok {
		return envVar{}, fmt.Errorf("missing DATABASE_URL")
	}
	ev.DBAddr = v

	return ev, nil
}
