package geo

import (
	"errors"
	"fmt"
	"net/url"
)

type (
	bingPoint struct {
		Coords []float64 `json:"coordinates"`
	}
	bingAddress struct {
		Address string `json:"formattedAddress"`
		Street  string `json:"addressLine"`
		Country string `json:"countryRegion"`
		Zip     string `json:"postalCode"`
		City    string `json:"locality"`
		State   string `json:"adminDistrict"`
	}
	bingResource struct {
		bingPoint   `json:"point"`
		bingAddress `json:"address"`
	}
	bingResourceSet struct {
		Total     int            `json:"estimatedTotal"`
		Resources []bingResource `json:"resources"`
	}
	bingResult struct {
		Set    []bingResourceSet `json:"resourceSets"`
		Status string            `json:"statusDescription"`
	}
	bingAPI struct {
		Geo, Img string
		Config
	}
)

func (api *bingAPI) Location(loc Location) (Result, error) {
	qry := url.Values{}
	qry.Add("o", "json")

	url := fmt.Sprintf("%s/%s?%s", api.Geo, loc, qry.Encode())
	return api.bingGeoService(url, loc.String())
}

func (api *bingAPI) Address(address string) (Result, error) {
	qry := url.Values{}
	qry.Add("q", address)
	qry.Add("o", "json")
	qry.Add("maxResults", "1")

	url := fmt.Sprintf("%s?%s", api.Geo, qry.Encode())
	return api.bingGeoService(url, address)
}

func (api *bingAPI) Image(address []string, opts MapOptions) ([]byte, error) {
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

func (api *bingAPI) bingGeoService(url, qry string) (res Result, err error) {
	var v bingResult
	if err = fetchJSON(api.Fetcher, url, &v); err != nil {
		return
	}

	if v.Status != "OK" {
		return res, errors.New(v.Status)
	} else if len(v.Set[0].Resources) == 0 {
		return res, errors.New("not found") //nothing found.
	}

	resx := v.Set[0].Resources[0]

	return Result{
		Query:    qry,
		Address:  resx.Address,
		Street:   resx.Street,
		Country:  resx.Country,
		Zip:      resx.Zip,
		City:     resx.City,
		State:    resx.State,
		Location: Location{resx.Coords[0], resx.Coords[1]},
	}, nil
}
