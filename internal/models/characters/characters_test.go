package characters

import (
	"context"
	"testing"

	"github.com/kagelui/marvel-forwarder/internal/testutil"
)

func TestCharacterSlice_saveWithTx(t *testing.T) {
	tests := []struct {
		name    string
		fixture CharacterSlice
		s       CharacterSlice
	}{
		{
			name: "nil",
			s:    nil,
		},
		{
			name: "some characters",
			s: []Character{
				{
					ID:          9312,
					Name:        "man",
					Description: "man0",
				},
				{
					ID:          831624,
					Name:        "woman",
					Description: "woman0",
				},
			},
		},
		{
			name: "with existing",
			fixture: []Character{
				{
					ID:          831624,
					Name:        "girl",
					Description: "girl0",
				},
			},
			s: []Character{
				{
					ID:          9312,
					Name:        "man",
					Description: "man0",
				},
				{
					ID:          831624,
					Name:        "woman",
					Description: "woman0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := db.MustBegin()
			testutil.Ok(t, tt.fixture.saveWithTx(context.TODO(), tx))
			testutil.Ok(t, tt.s.saveWithTx(context.TODO(), tx))
			for _, c := range tt.s {
				ch, e := GetCharacter(context.TODO(), tx, c.ID)
				testutil.Ok(t, e)
				testutil.Equals(t, c, ch)
			}
			testutil.Ok(t, tx.Rollback())
		})
	}
}
