package handler

import (
	"context"
	"fmt"
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
