package handler

import (
	"context"
	"net/http"

	"github.com/kagelui/marvel-forwarder/internal/pkg/web"
)

type characterStore interface {
	GetCharacterIDs(ctx context.Context) ([]int, error)
}

// GetMarvelCharacterList returns a list of marvel characters' ID
func GetMarvelCharacterList(s characterStore) web.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		ids, err := s.GetCharacterIDs(ctx)
		if err != nil {
			return web.NewError(err, "character ID list error")
		}
		web.RespondJSON(ctx, w, ids, nil)
		return nil
	}
}

// GetMarvelCharacterDetail returns a marvel characters' detail
func GetMarvelCharacterDetail(w http.ResponseWriter, r *http.Request) error {
	return nil
}
