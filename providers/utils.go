package providers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/muesli/cache2go"
)

// Accessing a new cache table for the first time will create it.
var cache = cache2go.Cache("gegeo")

func Cache(url, key string) ([]byte, error) {
	if len(key) > 0 {
		url = url + "&key=" + key
	}

	if res, err := cache.Value(url); err == nil {
		return res.Data().([]byte), nil
	}

	fmt.Println("retrieving:", url)
	resp, err := http.Get(url)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	cache.Add(url, 1*time.Minute, b)
	return b, err
}
