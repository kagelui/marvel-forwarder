package characters

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/kagelui/marvel-forwarder/internal/models/characters"
	"github.com/kagelui/marvel-forwarder/internal/pkg/web"
)

// ModelStore contains a reference to the DB connection and provides the service to handlers
type ModelStore struct {
	DB characters.Inquirer
}

// GetCharacterIDs returns the ID of all characters in the DB
func (m *ModelStore) GetCharacterIDs(ctx context.Context) ([]int, error) {
	chs, err := characters.GetCharacters(ctx, m.DB)
	if err != nil {
		return nil, err
	}

	result := make([]int, len(chs))
	for i, c := range chs {
		result[i] = c.ID
	}

	return result, err
}

// GetCharacter returns the character with the given id
func (m *ModelStore) GetCharacter(ctx context.Context, id int) (characters.Character, error) {
	ch, err := characters.GetCharacter(ctx, m.DB, id)
	switch {
	case err == sql.ErrNoRows:
		return characters.Character{}, web.Error{
			Status: http.StatusNotFound,
			Code:   "no_such_character",
			Desc:   "no such character",
		}
	case err != nil:
		return characters.Character{}, err
	}

	return ch, nil
}
