package providers

import (
	"io/ioutil"
	"net/http"
)

func GetBody(url, key string) ([]byte, error) {
	if len(key) > 0 {
		url = url + "&key=" + key
	}

	// fmt.Println("google:", url)
	resp, err := http.Get(url)

	// fmt.Println(resp)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
