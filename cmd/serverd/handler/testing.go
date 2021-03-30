package handler

import (
	"context"

	"github.com/kagelui/marvel-forwarder/internal/models/characters"
)

type mockStore struct {
	getCharacterIDsFn    func(context.Context) ([]int, error)
	getCharacterDetailFn func(ctx context.Context, id int) (characters.Character, error)
}

func (s mockStore) GetCharacterIDs(ctx context.Context) ([]int, error) {
	if s.getCharacterIDsFn != nil {
		return s.getCharacterIDsFn(ctx)
	}
	return nil, nil
}

func (s mockStore) GetCharacter(ctx context.Context, id int) (characters.Character, error) {
	if s.getCharacterDetailFn != nil {
		return s.getCharacterDetailFn(ctx, id)
	}
	return characters.Character{}, nil
}
