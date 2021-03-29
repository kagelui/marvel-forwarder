package characters

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Character contains the information of a character used in this app
type Character struct {
	ID          int    `db:"external_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

// CharacterSlice represents a slice of characters
type CharacterSlice []Character

// Save inserts the characters into DB, updating upon conflict of external_id
func (s CharacterSlice) Save(ctx context.Context, db *sqlx.DB) error {
	if len(s) == 0 {
		return nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	if err = s.saveWithTx(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (s CharacterSlice) saveWithTx(ctx context.Context, tx *sqlx.Tx) error {
	if len(s) == 0 {
		return nil
	}

	u := unique(s)
	insertQuery := `INSERT INTO characters (external_id, name, description) VALUES `
	positionStrSlice := make([]string, len(u))
	insertParams := make([]interface{}, 0)

	for i, c := range u {
		positionStrSlice[i] = fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3)
		insertParams = append(insertParams, c.ID, c.Name, c.Description)
	}

	insertQuery += strings.Join(positionStrSlice, ", ")
	insertQuery += ` ON CONFLICT (external_id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description`

	_, err := tx.ExecContext(ctx, insertQuery, insertParams...)
	return err
}

func unique(s []Character) []Character {
	m := make(map[int]Character)
	for _, c := range s {
		m[c.ID] = c
	}
	result := make([]Character, 0)
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

type inquirer interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// GetCharacter returns the Character struct for the given external ID
func GetCharacter(ctx context.Context, db inquirer, extID int) (Character, error) {
	var ch Character
	if err := db.GetContext(ctx, &ch, "SELECT external_id, name, description FROM characters WHERE external_id = $1", extID); err != nil {
		return Character{}, err
	}
	return ch, nil
}
