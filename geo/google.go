package geo

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type (
	googleGeometry struct {
		Location `json:"location"`
	}
	googleComponents struct {
		Long  string   `json:"long_name"`
		Short string   `json:"short_name"`
		Types []string `json:"types"`
	}
	googleResult struct {
		Compenents []googleComponents `json:"address_components"`
		Geometry   googleGeometry     `json:"geometry"`
		Address    string             `json:"formatted_address"`
	}
	googleResults struct {
		Results []googleResult `json:"results"`
		Status  string         `json:"status"`
	}
	googleAPI struct {
		Geo, Img string
		Config
	}
)

func (api *googleAPI) Location(loc Location) (Result, error) {
	qry := url.Values{}
	qry.Add("key", api.APIKey)
	qry.Add("latlng", loc.String())

	url := fmt.Sprintf("%s?%s", api.Geo, qry.Encode())
	return api.googleGeoService(url, loc.String())
}

func (api *googleAPI) Address(address string) (Result, error) {
	qry := url.Values{}
	qry.Add("key", api.APIKey)
	qry.Add("address", address)

	url := fmt.Sprintf("%s?%s", api.Geo, qry.Encode())
	return api.googleGeoService(url, address)
}

func (api *googleAPI) Image(address []string, opts MapOptions) ([]byte, error) {
	qry := url.Values{}
	qry.Add("key", api.APIKey)
	qry.Add("markers", strings.Join(address, "|"))
	qry.Add("size", opts.Size.String())

	if opts.Zoom > 0 {
		qry.Add("zoom", fmt.Sprintf("%v", opts.Zoom))
	}
	if opts.Scale > 0 {
		qry.Add("scale", fmt.Sprintf("%v", opts.Scale))
	}

	url := fmt.Sprintf("%s?%s", api.Img, qry.Encode())
	return fetch(api.Fetcher, url)
}

func (api *googleAPI) googleGeoService(url, qry string) (res Result, err error) {
	var result googleResults

	if err = fetchJSON(api.Fetcher, url, &result); err != nil {
		return
	}

	if result.Status != "OK" {
		return res, errors.New("result: " + result.Status)
	}

	return result.toGeoResult(qry), nil
}

func (r googleResults) toGeoResult(qry string) Result {
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

	// for key, val := range dic {
	// 	log.Println("-", key, "=>", val)
	// }

	return Result{
		Query:    qry,
		Street:   street,
		Country:  dic["country"],
		City:     dic["locality"],
		Zip:      dic["postal_code"],
		State:    state,
		Location: first.Geometry.Location,
		Address:  first.Address,
	}
}
