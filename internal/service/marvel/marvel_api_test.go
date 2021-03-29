package marvel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/kagelui/marvel-forwarder/internal/testutil"
)

func TestApiClient_requestHash(t *testing.T) {
	type fields struct {
		PublicKey  string
		PrivateKey string
	}
	type args struct {
		ts int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "valid",
			fields: fields{
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "265d12b39c12c21e267f5cc97137d5b0",
			},
			args: args{ts: 8312573},
			want: "faac4afc7298a4230a8594b0bcc680ae",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := ApiClient{
				PublicKey:  tt.fields.PublicKey,
				PrivateKey: tt.fields.PrivateKey,
			}
			got := ac.requestHash(tt.args.ts)
			testutil.Equals(t, tt.want, got)
		})
	}
}

type roundTripFunc func(req *http.Request) *http.Response

// RoundTrip implements RoundTripper
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// newTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestApiClient_retrieveOneBatchCharacters(t *testing.T) {
	type fields struct {
		Client     *http.Client
		PublicKey  string
		PrivateKey string
		Retries    int
	}
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    responseData
		wantErr string
	}{
		{
			name: "error",
			fields: fields{
				Client: newTestClient(func(req *http.Request) *http.Response {
					if !strings.Contains(req.URL.String(), "006127f9ec4cdd9da3973a1090fa1a75") {
						return &http.Response{StatusCode: http.StatusOK}
					}
					return &http.Response{StatusCode: http.StatusInternalServerError}
				}),
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "random",
				Retries:    0,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want:    responseData{},
			wantErr: "result error",
		},
		{
			name: "normal",
			fields: fields{
				Client: newTestClient(func(req *http.Request) *http.Response {
					if !strings.Contains(req.URL.String(), "006127f9ec4cdd9da3973a1090fa1a75") {
						return &http.Response{StatusCode: http.StatusBadRequest}
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(testutil.MustOpen("testdata/example.json")),
						Header:     make(http.Header),
					}
				}),
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "random",
				Retries:    0,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: responseData{
				Offset: 0,
				Limit:  10,
				Total:  1493,
				Count:  10,
				Results: []characterData{
					{
						ID:          1011334,
						Name:        "3-D Man",
						Description: "",
					},
				},
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := ApiClient{
				Client:     tt.fields.Client,
				PublicKey:  tt.fields.PublicKey,
				PrivateKey: tt.fields.PrivateKey,
			}
			got, err := ac.retrieveOneBatchCharacters(context.TODO(), tt.args.offset, tt.args.limit)
			testutil.CompareError(t, tt.wantErr, err)
			if err == nil {
				t.Log(got)
				testutil.Equals(t, tt.want, got)
			}
		})
	}
}

func TestApiClient_RetrieveCharacters(t *testing.T) {
	alwaysNaughtyClient := newTestClient(func(req *http.Request) *http.Response {
		if !strings.Contains(req.URL.String(), "006127f9ec4cdd9da3973a1090fa1a75") {
			return &http.Response{StatusCode: http.StatusOK}
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}
	})
	occasionallyNaughtyClient := newTestClient(func(req *http.Request) *http.Response {
		if !strings.Contains(req.URL.String(), "006127f9ec4cdd9da3973a1090fa1a75") {
			return &http.Response{StatusCode: http.StatusOK}
		}
		offsets, ok := req.URL.Query()["offset"]
		if !ok || len(offsets) < 1 {
			return &http.Response{StatusCode: http.StatusBadRequest}
		}
		if offsets[0] == "300" {
			return &http.Response{StatusCode: http.StatusInternalServerError}
		}
		fileName := fmt.Sprintf("testdata/%v.json", offsets[0])
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(testutil.MustOpen(fileName)),
			Header:     make(http.Header),
		}
	})
	normalClient := newTestClient(func(req *http.Request) *http.Response {
		if !strings.Contains(req.URL.String(), "006127f9ec4cdd9da3973a1090fa1a75") {
			return &http.Response{StatusCode: http.StatusBadRequest}
		}
		offsets, ok := req.URL.Query()["offset"]
		if !ok || len(offsets) < 1 {
			return &http.Response{StatusCode: http.StatusBadRequest}
		}
		fileName := fmt.Sprintf("testdata/%v.json", offsets[0])
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(testutil.MustOpen(fileName)),
			Header:     make(http.Header),
		}
	})

	normalResult := make([]Character, 0)
	bytes, _ := ioutil.ReadAll(testutil.MustOpen("testdata/result.json"))
	_ = json.Unmarshal(bytes, &normalResult)

	type fields struct {
		Client     *http.Client
		APIAddr    string
		PublicKey  string
		PrivateKey string
		Retries    int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Character
		wantErr string
	}{
		{
			name: "normal",
			fields: fields{
				Client:     normalClient,
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "--",
				Retries:    0,
			},
			want:    normalResult,
			wantErr: "",
		},
		{
			name: "error in first call",
			fields: fields{
				Client:     alwaysNaughtyClient,
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "--",
				Retries:    0,
			},
			want:    nil,
			wantErr: "result error",
		},
		{
			name: "error in parallel call",
			fields: fields{
				Client:     occasionallyNaughtyClient,
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "--",
				Retries:    0,
			},
			want:    nil,
			wantErr: "error in parallel run",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := ApiClient{
				Client:     tt.fields.Client,
				APIAddr:    tt.fields.APIAddr,
				PublicKey:  tt.fields.PublicKey,
				PrivateKey: tt.fields.PrivateKey,
				Retries:    tt.fields.Retries,
			}
			got, err := ac.RetrieveCharacters(context.TODO())
			testutil.CompareError(t, tt.wantErr, err)
			if err == nil {
				testutil.Equals(t, sortCharacters(tt.want), sortCharacters(got))
			}
		})
	}
}

func sortCharacters(characters []Character) []Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].ID < characters[j].ID
	})
	return characters
}
