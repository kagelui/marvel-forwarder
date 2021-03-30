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
		want    CharacterSlice
	}{
		{
			name: "nil",
			s:    nil,
			want: nil,
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
			want: []Character{
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
			want: []Character{
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
			name: "repeated values with existing",
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
				{
					ID:          831624,
					Name:        "woman",
					Description: "woman1",
				},
			},
			want: []Character{
				{
					ID:          9312,
					Name:        "man",
					Description: "man0",
				},
				{
					ID:          831624,
					Name:        "woman",
					Description: "woman1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := db.MustBegin()
			testutil.Ok(t, tt.fixture.SaveWithTx(context.TODO(), tx))
			testutil.Ok(t, tt.s.SaveWithTx(context.TODO(), tx))
			for _, c := range tt.want {
				ch, e := GetCharacter(context.TODO(), tx, c.ID)
				testutil.Ok(t, e)
				testutil.Equals(t, c, ch)
			}
			testutil.Ok(t, tx.Rollback())
		})
	}
}

func TestGetCharacters(t *testing.T) {
	tests := []struct {
		name    string
		fixture CharacterSlice
		want    []Character
		wantErr string
	}{
		{
			name:    "return empty slice when DB is empty",
			fixture: nil,
			want:    []Character{},
			wantErr: "",
		},
		{
			name: "some records",
			fixture: []Character{
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
			},
			want: []Character{
				{
					ID:          186824,
					Name:        "Stick",
					Description: "teacher to some broke lawyer",
				},
				{
					ID:          941356,
					Name:        "Daredevil",
					Description: "some broke lawyer",
				},
				{
					ID:          1082344,
					Name:        "Kingpin",
					Description: "some bulky and rich villain",
				},
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := db.MustBegin()
			tx.MustExec(`TRUNCATE characters`)
			testutil.Ok(t, tt.fixture.SaveWithTx(context.TODO(), tx))
			got, err := GetCharacters(context.TODO(), tx)
			testutil.CompareError(t, tt.wantErr, err)
			if err == nil {
				testutil.Equals(t, tt.want, got)
			}
			testutil.Ok(t, tx.Rollback())
		})
	}
}
