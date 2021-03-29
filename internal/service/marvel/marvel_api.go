package marvel

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/kagelui/marvel-forwarder/internal/models/characters"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/kagelui/marvel-forwarder/internal/pkg/loglib"
)

const (
	apiLimit = 100
)

type ApiClient struct {
	Client     *http.Client
	PublicKey  string
	PrivateKey string
	APIAddr    string
	Retries    int
}

// RetrieveCharacters retrieves all the characters from the API
func (ac ApiClient) RetrieveCharacters(ctx context.Context) (characters.CharacterSlice, error) {
	lg := loglib.GetLogger(ctx)
	result := make([]characters.Character, 0)
	firstTrunk, err := ac.retrieveOneBatchCharacters(ctx, 0, apiLimit)
	if err != nil {
		return nil, err
	}
	result = append(result, responseToCharacters(firstTrunk)...)

	total := firstTrunk.Total
	numConnections := total / apiLimit
	var wg sync.WaitGroup
	var hasError bool
	for i := 1; i <= numConnections; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			lg.InfoF("Starting runner %d", num)
			resp, err := ac.retrieveOneBatchCharacters(ctx, num*apiLimit, apiLimit)
			if err != nil {
				lg.ErrorF(err.Error())
				hasError = true
			}
			result = append(result, responseToCharacters(resp)...)
			lg.InfoF("End runner %d", num)
		}(i)
	}
	wg.Wait()

	if hasError {
		return nil, fmt.Errorf("error in parallel run")
	}

	return result, nil
}

type response struct {
	Code            int          `json:"code"`
	Status          string       `json:"status"`
	Copyright       string       `json:"copyright"`
	AttributionText string       `json:"attributionText"`
	AttributionHTML string       `json:"attributionHTML"`
	Etag            string       `json:"etag"`
	Data            responseData `json:"data"`
}

type responseData struct {
	Offset  int             `json:"offset"`
	Limit   int             `json:"limit"`
	Total   int             `json:"total"`
	Count   int             `json:"count"`
	Results []characterData `json:"results"`
}

type characterData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ac ApiClient) retrieveOneBatchCharacters(ctx context.Context, offset, limit int) (responseData, error) {
	ts := time.Now().Unix()
	hash := ac.requestHash(ts)
	addr := fmt.Sprintf("%s?ts=%v&apikey=%s&hash=%s&offset=%d&limit=%d", ac.APIAddr, ts, ac.PublicKey, hash, offset, limit)
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return responseData{}, err
	}

	var resp *http.Response
	var e error
	if err := withRetries(ctx, func() error {
		resp, e = ac.Client.Do(req)
		if e != nil {
			return e
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("result error")
		}
		return nil
	}, ac.Retries); err != nil {
		return responseData{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseData{}, err
	}
	resp.Body.Close()

	var r response
	if er := json.Unmarshal(data, &r); er != nil {
		return responseData{}, er
	}

	return r.Data, nil
}

func (ac ApiClient) requestHash(ts int64) string {
	h := md5.New()
	_, _ = io.WriteString(h, fmt.Sprintf("%v%v%v", ts, ac.PrivateKey, ac.PublicKey))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func withRetries(ctx context.Context, callback func() error, retries int) error {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 5 * time.Second
	b.RandomizationFactor = 0
	b.MaxElapsedTime = time.Minute
	bo := backoff.WithContext(backoff.WithMaxRetries(b, uint64(retries)), ctx)

	return backoff.Retry(callback, bo)
}

func responseToCharacters(rd responseData) []characters.Character {
	result := make([]characters.Character, len(rd.Results))
	for i, one := range rd.Results {
		result[i] = characters.Character{
			ID:          one.ID,
			Name:        one.Name,
			Description: one.Description,
		}
	}
	return result
}
