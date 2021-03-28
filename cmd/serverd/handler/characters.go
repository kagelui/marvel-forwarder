package handler

import (
	"net/http"

	"github.com/kagelui/marvel-forwarder/internal/pkg/web"
)

// GetMarvelCharacterList returns a list of marvel characters' ID
func GetMarvelCharacterList(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	web.RespondJSON(ctx, w, []int{31264, 12983}, nil)
	return nil
}

// GetMarvelCharacterDetail returns a marvel characters' detail
func GetMarvelCharacterDetail(w http.ResponseWriter, r *http.Request) error {
	return nil
}
