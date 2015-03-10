package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/harboe/gogeo/providers"
)

const (
	geourl = "https://maps.googleapis.com/maps/api/geocode/json"
	imgurl = "https://maps.googleapis.com/maps/api/staticmap"
)

type (
	geometry struct {
		providers.Location `json:"location"`
	}
	result struct {
		geometry `json:"geometry"`
		Address  string `json:"formatted_address"`
	}
	results struct {
		Results []result `json:"results"`
		Status  string   `json:"status"`
	}
	geoService struct {
		key string
	}
)

func init() {
	providers.RegisterGeo("google", func(key string) (providers.GeoService, error) {
		return geoService{key}, nil
	})
}

func (g geoService) Location(loc providers.Location) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("latlng", loc.String())

	url := fmt.Sprintf("%s?%s", geourl, qry.Encode())
	return g.googleGeoService(url)
}

func (g geoService) Address(address string) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("address", address)

	url := fmt.Sprintf("%s?%s", geourl, qry.Encode())
	return g.googleGeoService(url)
}

func (g geoService) Static(address []string, opts providers.MapOptions) ([]byte, error) {
	qry := url.Values{}

	markers := ""
	for _, a := range address {
		markers += fmt.Sprintf("%s|", a)
	}
	qry.Add("markers", markers)

	if opts.Zoom > 0 {
		qry.Add("zoom", fmt.Sprintf("%v", opts.Zoom))
	}
	if opts.Scale > 0 {
		qry.Add("scale", fmt.Sprintf("%v", opts.Scale))
	}
	qry.Add("size", opts.Size.String())

	url := fmt.Sprintf("%s?%s", imgurl, qry.Encode())
	return g.getResponseBody(url)
}

func (g geoService) googleGeoService(url string) (providers.Result, error) {
	b, err := g.getResponseBody(url)

	if err != nil {
		return providers.Result{}, err
	}

	var result results

	if err := json.Unmarshal(b, &result); err != nil {
		return providers.Result{}, err
	}

	if result.Status != "OK" {
		return providers.Result{}, errors.New("result: " + result.Status)
	}

	return providers.Result{
		Location: result.Results[0].Location,
		Address:  result.Results[0].Address,
	}, nil
}

func (g geoService) getResponseBody(url string) ([]byte, error) {
	if len(g.key) > 0 {
		url = url + "&key=" + g.key
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
