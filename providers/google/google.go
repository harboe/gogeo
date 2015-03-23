package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	components struct {
		Long  string   `json:"long_name"`
		Short string   `json:"short_name"`
		Types []string `json:"types"`
	}
	result struct {
		Compenents []components `json:"address_components"`
		geometry   `json:"geometry"`
		Address    string `json:"formatted_address"`
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
	providers.Register("google", func(key string) (providers.GeoService, error) {
		return geoService{key}, nil
	})
}

func (g geoService) Location(loc providers.Location) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("latlng", loc.String())

	url := fmt.Sprintf("%s?%s", geourl, qry.Encode())
	return g.googleGeoService(url, loc.String())
}

func (g geoService) Address(address string) (providers.Result, error) {
	qry := url.Values{}
	qry.Add("address", address)

	url := fmt.Sprintf("%s?%s", geourl, qry.Encode())
	return g.googleGeoService(url, address)
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
	return providers.Cache(url, g.key)
}

func (g geoService) googleGeoService(url, qry string) (providers.Result, error) {
	b, err := providers.Cache(url, g.key)

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

	// fmt.Printf("%#v\n", result)
	// fmt.Println(string(b))

	return result.toProviderResult(qry), nil
}

func (r results) toProviderResult(qry string) providers.Result {
	first := r.Results[0]
	dic := map[string]string{}

	for _, c := range first.Compenents {
		dic[c.Types[0]] = c.Long
	}

	state := ""
	if val, ok := dic["administrative_area_level_1"]; ok {
		state = val
	}
	street := ""
	if val, ok := dic["street_number"]; ok {
		street = fmt.Sprintf("%s %v", dic["route"], val)
	} else {
		street = dic["route"]
	}

	for key, val := range dic {
		log.Println("-", key, "=>", val)
	}

	return providers.Result{
		Query:    qry,
		Street:   street,
		Country:  dic["country"],
		City:     dic["locality"],
		Zip:      dic["postal_code"],
		State:    state,
		Location: first.Location,
		Address:  first.Address,
	}
}
