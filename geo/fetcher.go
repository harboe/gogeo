package geo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	Fetcher interface {
		Do(r *http.Request) (*http.Response, error)
	}
	FetcherFunc func(*http.Request) (*http.Response, error)
)

func (fn FetcherFunc) Do(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func fetch(c Fetcher, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return []byte{}, err
	}

	res, err := c.Do(req)

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func fetchJSON(c Fetcher, url string, v interface{}) error {
	b, err := fetch(c, url)

	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
