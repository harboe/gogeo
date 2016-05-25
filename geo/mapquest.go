package geo

import (
	"errors"
	"fmt"
	"net/url"
)

type (
	mqLocation struct {
		Location `json:"latLng"`
		Street   string `json:"street"`
		Country  string `json:"adminArea1"`
		Zip      string `json:"postalCode"`
		City     string `json:"admininArea5"`
		State    string `json:"adminArea3"`
	}
	mqResult struct {
		Location []mqLocation `json:"locations"`
	}
	mqInfo struct {
		Status  int      `json:"statuscode"`
		Message []string `json:"messages"`
	}
	mqPayload struct {
		Results []mqResult `json:"results"`
		Info    mqInfo     `json:"info"`
	}
	mapquestAPI struct {
		Geo, Img string
		Config
	}
)

func (api *mapquestAPI) Location(loc Location) (Result, error) {
	qry := url.Values{}
	qry.Add("location", loc.String())

	url := fmt.Sprintf("%s%s?%s", api.Geo, "reverse", qry.Encode())
	return api.toProviderResult(url, loc.String())
}

func (api *mapquestAPI) Address(address string) (Result, error) {
	qry := url.Values{}
	qry.Add("location", address)
	qry.Add("maxResults", "1")
	qry.Add("thumbMaps", "false")

	url := fmt.Sprintf("%s%s?%s", api.Geo, "address", qry.Encode())
	return api.toProviderResult(url, address)
}

func (api *mapquestAPI) Image(markers []string, options MapOptions) ([]byte, error) {
	return []byte{}, fmt.Errorf("not implemeted")
}

func (api *mapquestAPI) toProviderResult(url, qry string) (res Result, err error) {
	var p mqPayload
	if err = fetchJSON(api.Fetcher, url, &p); err != nil {
		return
	}

	if len(p.Results[0].Location) == 0 {
		return res, errors.New("not found")
	}

	l := p.Results[0].Location[0]
	// addr := "muuha..."
	return Result{
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
