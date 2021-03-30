package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kagelui/marvel-forwarder/internal/models/characters"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kagelui/marvel-forwarder/internal/pkg/web"
	"github.com/kagelui/marvel-forwarder/internal/testutil"
)

func TestGetMarvelCharacterList(t *testing.T) {
	type args struct {
		s characterStore
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
		expectedBody string
	}{
		{
			name: "naughty store",
			args: args{s: mockStore{
				getCharacterIDsFn: func(ctx context.Context) ([]int, error) {
					return []int{}, fmt.Errorf("mock error")
				},
			}},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"internal_error","error_description":"Sorry, there was a problem. Please try again later."}`,
		},
		{
			name: "normal",
			args: args{s: mockStore{
				getCharacterIDsFn: func(ctx context.Context) ([]int, error) {
					return []int{391264, 831256}, nil
				},
			}},
			expectedCode: http.StatusOK,
			expectedBody: `[391264,831256]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			web.Handler{H: GetMarvelCharacterList(tt.args.s)}.ServeHTTP(rr, req)
			testutil.Equals(t, tt.expectedCode, rr.Code)
			testutil.Equals(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetMarvelCharacterDetail(t *testing.T) {
	type args struct {
		s characterStore
	}
	tests := []struct {
		name         string
		args         args
		id           string
		expectedCode int
		expectedBody string
	}{
		{
			name: "not found",
			args: args{s: mockStore{getCharacterDetailFn: func(ctx context.Context, id int) (characters.Character, error) {
				if id != 83253 {
					return characters.Character{}, nil
				}
				return characters.Character{}, &web.Error{
					Status: http.StatusNotFound,
					Code:   "no_such_character",
					Desc:   "no such character",
				}
			}}},
			id:           "83253",
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"no_such_character","error_description":"no such character"}`,
		},
		{
			name: "wonky store",
			args: args{s: mockStore{getCharacterDetailFn: func(ctx context.Context, id int) (characters.Character, error) {
				if id != 28663 {
					return characters.Character{}, nil
				}
				return characters.Character{}, fmt.Errorf("mock error")
			}}},
			id:           "28663",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"internal_error","error_description":"Sorry, there was a problem. Please try again later."}`,
		},
		{
			name:         "malformed ID",
			args:         args{s: mockStore{}},
			id:           "not a number",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"malformed_id","error_description":"malformed ID"}`,
		},
		{
			name: "nice and peaceful",
			args: args{s: mockStore{getCharacterDetailFn: func(ctx context.Context, id int) (characters.Character, error) {
				if id != 832634 {
					return characters.Character{}, fmt.Errorf("data error")
				}
				return characters.Character{ID: 832654, Name: "Daredevil", Description: "some broke lawyer"}, nil
			}}},
			id:           "832634",
			expectedCode: http.StatusOK,
			expectedBody: `{"ID":832654,"Name":"Daredevil","Description":"some broke lawyer"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()
			web.Handler{H: GetMarvelCharacterDetail(tt.args.s)}.ServeHTTP(rr, req)
			testutil.Equals(t, tt.expectedCode, rr.Code)
			testutil.Equals(t, tt.expectedBody, rr.Body.String())
		})
	}
}
