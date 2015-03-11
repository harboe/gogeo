package bing

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
	geourl = "https://dev.virtualearth.net/REST/v1/Locations"
	// imgurl  = "https://maps.googleapis.com/maps/api/staticmap"
	// locurl = "https://dev.virtualearth.net/REST/v1/Locations"
)

type (
	point struct {
		// Type   string    `json:"typ"`
		Coords []float64 `json:"coordinates"`
	}
	address struct {
		Address string `json:"formattedAddress"`
	}
	resource struct {
		point   `json:"point"`
		address `json:"address"`
	}
	resourceSet struct {
		Total     int        `json:"estimatedTotal"`
		Resources []resource `json:"resources"`
	}
	result struct {
		Set    []resourceSet `json:"resourceSets"`
		Status string        `json:"statusDescription"`
	}

	bingService struct {
		key string
	}
)

func init() {
	providers.RegisterGeo("bing", func(key string) (providers.GeoService, error) {
		// validate key length, since bing requires a key
		return bingService{key}, nil
	})
}

func (b bingService) Location(loc providers.Location) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("o", "json")

	url := fmt.Sprintf("%s/%s?%s", geourl, loc, qry.Encode())
	return b.bingGeoService(url, loc.String())
}

func (b bingService) Address(address string) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("q", address)
	qry.Add("o", "json")
	qry.Add("maxResults", "1")

	url := fmt.Sprintf("%s?%s", geourl, qry.Encode())
	return b.bingGeoService(url, address)
}

func (b bingService) Static(address []string, opts providers.MapOptions) ([]byte, error) {
	// qry := url.Values{}

	// markers := ""
	// for _, a := range address {
	// 	markers += fmt.Sprintf("%s|", a)
	// }
	// qry.Add("markers", markers)

	// if opts.Zoom > 0 {
	// 	qry.Add("zoom", fmt.Sprintf("%v", opts.Zoom))
	// }
	// if opts.Scale > 0 {
	// 	qry.Add("scale", fmt.Sprintf("%v", opts.Scale))
	// }
	// qry.Add("size", opts.Size.String())

	// url := fmt.Sprintf("%s?%s", imgurl, qry.Encode())
	// return g.getResponseBody(url)
	return []byte{}, errors.New("not implemeted")
}

func (b bingService) bingGeoService(url, qry string) (providers.Result, error) {
	body, err := b.getResponseBody(url)

	if err != nil {
		return providers.Result{}, err
	}

	var v result

	if err := json.Unmarshal(body, &v); err != nil {
		return providers.Result{}, err
	}

	if v.Status != "OK" {
		return providers.Result{}, errors.New(v.Status)
	} else if len(v.Set[0].Resources) == 0 {
		return providers.Result{}, errors.New("not found") //nothing found.
	}

	resx := v.Set[0].Resources[0]

	return providers.Result{
		Query:    qry,
		Address:  resx.Address,
		Location: providers.Location{resx.Coords[0], resx.Coords[1]},
	}, nil
}

func (b bingService) getResponseBody(url string) ([]byte, error) {
	if len(b.key) > 0 {
		url = url + "&key=" + b.key
	}

	resp, err := http.Get(url)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
