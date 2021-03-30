package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kagelui/marvel-forwarder/internal/models/characters"
	"github.com/kagelui/marvel-forwarder/internal/pkg/web"
)

type characterStore interface {
	GetCharacterIDs(ctx context.Context) ([]int, error)
	GetCharacter(ctx context.Context, id int) (characters.Character, error)
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
func GetMarvelCharacterDetail(s characterStore) web.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			return &web.Error{
				Status: http.StatusBadRequest,
				Code:   "malformed_id",
				Desc:   "malformed ID",
				Err:    nil,
			}
		}

		info, err := s.GetCharacter(ctx, id)
		if err != nil {
			return web.NewError(err, "no such character")
		}
		web.RespondJSON(ctx, w, info, nil)
		return nil
	}
}
