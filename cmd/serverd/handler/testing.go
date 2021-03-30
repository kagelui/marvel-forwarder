package handler

import "context"

type mockStore struct {
	getCharacterIDsFn func(context.Context) ([]int, error)
}

func (s mockStore) GetCharacterIDs(ctx context.Context) ([]int, error) {
	if s.getCharacterIDsFn != nil {
		return s.getCharacterIDsFn(ctx)
	}
	return nil, nil
}
