package characters

import (
	"context"
	"testing"

	"github.com/kagelui/marvel-forwarder/internal/models/characters"
	"github.com/kagelui/marvel-forwarder/internal/testutil"
)

func TestModelStore_GetCharacter(t *testing.T) {
	type fixture struct {
		characters characters.CharacterSlice
	}
	tests := []struct {
		name    string
		f       fixture
		id      int
		want    characters.Character
		wantErr string
	}{
		{
			name:    "not found",
			f:       fixture{},
			id:      3182643,
			want:    characters.Character{},
			wantErr: "no such character",
		},
		{
			name: "not found",
			f: fixture{characters: []characters.Character{
				{
					ID:          941356,
					Name:        "Daredevil",
					Description: "some broke lawyer",
				},
				{
					ID:          186824,
					Name:        "Stick",
					Description: "teacher to some broke lawyer",
				},
				{
					ID:          1082344,
					Name:        "Kingpin",
					Description: "some bulky and rich villain",
				},
			}},
			id: 186824,
			want: characters.Character{
				ID:          186824,
				Name:        "Stick",
				Description: "teacher to some broke lawyer",
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := db.MustBegin()
			tx.MustExec(`TRUNCATE characters`)
			testutil.Ok(t, tt.f.characters.SaveWithTx(context.TODO(), tx))

			m := &ModelStore{
				DB: tx,
			}
			got, err := m.GetCharacter(context.TODO(), tt.id)
			testutil.CompareError(t, tt.wantErr, err)
			if err == nil {
				testutil.Equals(t, tt.want, got)
			}
			testutil.Ok(t, tx.Rollback())
		})
	}
}

func TestModelStore_GetCharacterIDs(t *testing.T) {
	type fixture struct {
		characters characters.CharacterSlice
	}
	tests := []struct {
		name    string
		f       fixture
		want    []int
		wantErr string
	}{
		{
			name:    "should return empty slice should DB be empty",
			f:       fixture{},
			want:    []int{},
			wantErr: "",
		},
		{
			name: "should return ID slice in asc order",
			f: fixture{characters: []characters.Character{
				{
					ID:          941356,
					Name:        "Daredevil",
					Description: "some broke lawyer",
				},
				{
					ID:          186824,
					Name:        "Stick",
					Description: "teacher to some broke lawyer",
				},
				{
					ID:          1082344,
					Name:        "Kingpin",
					Description: "some bulky and rich villain",
				},
			}},
			want:    []int{186824, 941356, 1082344},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := db.MustBegin()
			tx.MustExec(`TRUNCATE characters`)
			testutil.Ok(t, tt.f.characters.SaveWithTx(context.TODO(), tx))

			m := &ModelStore{
				DB: tx,
			}
			got, err := m.GetCharacterIDs(context.TODO())
			testutil.CompareError(t, tt.wantErr, err)
			if err == nil {
				testutil.Equals(t, tt.want, got)
			}
			testutil.Ok(t, tx.Rollback())
		})
	}
}
