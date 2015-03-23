package mapquest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/harboe/gogeo/providers"
)

const geourl = "https://open.mapquestapi.com/geocoding/v1/"

type (
	location struct {
		providers.Location `json:"latLng"`
		Street             string `json:"street"`
		Country            string `json:"adminArea1"`
		Zip                string `json:"postalCode"`
		City               string `json:"adminArea5"`
		State              string `json:"adminArea3"`
	}
	result struct {
		Location []location `json:"locations"`
	}
	info struct {
		Status  int      `json:"statuscode"`
		Message []string `json:"messages"`
	}
	payload struct {
		Results []result `json:"results"`
		Info    info     `json:"info"`
	}
	mqService struct {
		key string
	}
)

func init() {
	providers.Register("mapquest", func(key string) (providers.GeoService, error) {
		return mqService{key}, nil
	})
}

func (mq mqService) Location(loc providers.Location) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("location", loc.String())

	url := fmt.Sprintf("%s%s?%s", geourl, "reverse", qry.Encode())
	return mq.toProviderResult(url, loc.String())
}

func (mq mqService) Address(address string) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("location", address)
	qry.Add("maxResults", "1")
	qry.Add("thumbMaps", "false")

	url := fmt.Sprintf("%s%s?%s", geourl, "address", qry.Encode())
	return mq.toProviderResult(url, address)
}

func (mq mqService) Static(markers []string, options providers.MapOptions) ([]byte, error) {
	return []byte{}, nil
}

func (mq mqService) toProviderResult(url, qry string) (providers.Result, error) {
	b, err := providers.Cache(url, mq.key)

	if err != nil {
		return providers.Result{}, err
	}

	var p payload
	if err := json.Unmarshal(b, &p); err != nil {
		return providers.Result{}, err
	}

	if len(p.Results[0].Location) == 0 {
		return providers.Result{}, errors.New("not found")
	}

	l := p.Results[0].Location[0]
	// addr := "muuha..."
	return providers.Result{
		Query: qry,
		// Address:  addr,
		Street:   l.Street,
		Country:  l.Country,
		Zip:      l.Zip,
		City:     l.City,
		State:    l.State,
		Location: l.Location,
	}, nil

}
