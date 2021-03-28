package marvel

import (
	"context"
	"io/ioutil"
	"net/http"
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
				Retries:    1,
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
				Retries:    1,
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
				Client:     http.DefaultClient,
				PublicKey:  "006127f9ec4cdd9da3973a1090fa1a75",
				PrivateKey: "--",
				Retries:    1,
			},
			want:    nil,
			wantErr: "",
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
				testutil.Equals(t, tt.want, got)
			}
		})
	}
}
