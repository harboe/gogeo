package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type (
	Result struct {
		Address   string  `json:"address"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
	}
	GoogleLocation struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
	}
	GoogleGeometry struct {
		Location GoogleLocation `json:"location"`
	}
	GoogleResult struct {
		GoogleGeometry `json:"geometry"`
		Address        string `json:"formatted_address"`
	}
	GoogleResults struct {
		Results []GoogleResult `json:"results"`
		Status  string         `json:"status"`
	}
)

const (
	geourl = "https://maps.googleapis.com/maps/api/geocode/json"
	imgurl = "https://maps.googleapis.com/maps/api/staticmap"
)

func main() {
	http.HandleFunc("/", usageHandler)
	http.HandleFunc("/v1/geo", geoHandler)
	http.HandleFunc("/v1/geo.png", imgHandler)
	log.Println(http.ListenAndServe("localhost:8080", nil))
}

func imgHandler(w http.ResponseWriter, req *http.Request) {
	results, err := getResults(req)

	if len(results) == 0 || err != nil {
		displayError(w, req, err)
	} else {

		qry := url.Values{}
		markers := ""

		urlQry := req.URL.Query()
		size := urlQry.Get("size")

		if len(size) == 0 {
			size = "250x250"
		} else if strings.Index(size, "x") == -1 {
			size = size + "x" + size
		}

		qry.Add("size", size)

		if z := urlQry.Get("zoom"); len(z) > 0 {
			qry.Add("zoom", z)
		}

		if s := urlQry.Get("scale"); len(s) > 0 {
			qry.Add("scale", s)
		}

		for _, r := range results {
			markers += fmt.Sprintf("%v,%v|", r.Latitude, r.Longitude)
		}
		qry.Add("markers", markers)

		url := fmt.Sprintf("%s?%s", imgurl, qry.Encode())
		fmt.Println("url:", url)

		resp, err := http.Get(url)

		if err != nil {
			displayError(w, req, err)
			return
		}

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			displayError(w, req, err)
			return
		}

		w.Write(b)
	}
}

func geoHandler(w http.ResponseWriter, req *http.Request) {
	results, err := getResults(req)

	if len(results) == 0 || err != nil {
		displayError(w, req, err)
	} else {
		displayResult(w, req, results)
	}
}

func usageHandler(w http.ResponseWriter, req *http.Request) {
	usage := `
		GOGEO:
		------

		/v1/geo

		parameters:
		* address - The street address that you want to geocode, in the format used by the national postal service of the country concerned. 
		* lat - The textual latitude value for which you wish to obtain the closest, human-readable address.
		* lng - The textual longitude value for which you wish to obtain the closest, human-readable address.
		
		/v1/geo.png

		parameters:
		* address - The street address that you want to geocode, in the format used by the national postal service of the country concerned. 
		* lat - The textual latitude value for which you wish to obtain the closest, human-readable address.
		* lng - The textual longitude value for which you wish to obtain the closest, human-readable address.
		* size (optional)
		* zoom (optional)
		* scale (optional)

		`

	w.Write([]byte(usage))
}

func displayError(w http.ResponseWriter, req *http.Request, err error) {
	usageHandler(w, req)
	fmt.Fprintln(w, "error processing:", req.URL.RawQuery, "\n\t\terr:", err)
}

func displayResult(w http.ResponseWriter, req *http.Request, results []Result) {
	if b, err := json.Marshal(&results); err != nil {
		displayError(w, req, err)
	} else {
		w.Write(b)
	}
}

func getResults(req *http.Request) ([]Result, error) {
	url := req.URL
	values := url.Query()

	var results = []Result{}
	var err error

	if address, ok := values["address"]; ok {
		results, err = lookupAddress(address)
	}

	lat, hasLat := values["lat"]
	lng, hasLng := values["lng"]

	if hasLat && hasLng {
		results, err = lookupLocation(lat, lng)
	}

	return results, err
}

func lookupAddress(address []string) ([]Result, error) {
	results := []Result{}

	for _, a := range address {
		qry := url.Values{}
		qry.Add("address", a)

		url := fmt.Sprintf("%s?%s", geourl, qry.Encode())

		if r, err := lookup(url); err != nil {
			return results, err
		} else {
			results = append(results, r)
		}
	}

	return results, nil
}

func lookupLocation(lat, lng []string) ([]Result, error) {
	size := len(lat)

	if size > len(lng) {
		size = len(lng)
	}

	results := []Result{}

	for i := 0; i < size; i++ {
		qry := url.Values{}
		qry.Add("latlng", lat[i]+","+lng[i])

		url := fmt.Sprintf("%s?%s", geourl, qry.Encode())

		if r, err := lookup(url); err != nil {
			return results, err
		} else {
			results = append(results, r)
		}
	}

	return results, nil
}

func lookup(url string) (Result, error) {
	log.Println("url:", url)
	resp, err := http.Get(url)

	if err != nil {
		return Result{}, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	var result GoogleResults

	if err := json.Unmarshal(b, &result); err != nil {
		return Result{}, err
	}

	if result.Status != "OK" {
		log.Println("result:", result)
		return Result{}, nil
	}

	return Result{
		Latitude:  result.Results[0].Location.Latitude,
		Longitude: result.Results[0].Location.Longitude,
		Address:   result.Results[0].Address,
	}, nil
}
