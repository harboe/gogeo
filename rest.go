package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/harboe/gogeo/providers"
)

func RestService(port string) {
	router := httprouter.New()
	router.GET("/v1/:name/png", imgHandler)
	router.GET("/v1/:name/json", geoHandler)
	router.GET("/v1/:name/xml", geoHandler)
	router.GET("/v1/:name/txt", geoHandler)
	router.NotFound = helpHandler

	log.Fatal(http.ListenAndServe(port, router))
}

func geoHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	geo, err := restGeoService(ps.ByName("name"))

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var results providers.Results

	for _, l := range location(req) {
		if r, err := geo.Location(l); err == nil {
			results = append(results, r)
		}
	}

	for _, a := range address(req) {
		if r, err := geo.Address(a); err == nil {
			results = append(results, r)
		}
	}

	_, pretty := req.URL.Query()["pretty"]
	format := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	b, err := Marshal(format, &results, pretty)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func imgHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	geo, err := restGeoService(ps.ByName("name"))

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	opts := mapOptions(req)
	markers := address(req)

	for _, loc := range location(req) {
		markers = append(markers, loc.String())
	}

	b, err := geo.Static(markers, opts)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func helpHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("documentation not implemented"))
}

func restGeoService(name string) (providers.GeoService, error) {
	keyP := providerKeys[name]
	keyV := ""

	if keyP != nil {
		keyV = *keyP
	}

	return providers.Geo(name, keyV)
}

func mapOptions(req *http.Request) providers.MapOptions {
	qry := req.URL.Query()
	opts := providers.MapOptions{
		Size: providers.ParseSize(qry.Get("size")),
	}

	if z := qry.Get("zoom"); len(z) > 0 {
		opts.Zoom, _ = strconv.ParseUint(z, 0, 10)
	}

	if s := qry.Get("scale"); len(s) > 0 {
		opts.Scale, _ = strconv.ParseUint(s, 0, 10)
	}

	return opts
}

func location(req *http.Request) []providers.Location {
	lngs := req.URL.Query()["loc"]
	list := make([]providers.Location, len(lngs))

	for i := 0; i < len(lngs); i++ {
		list[i] = providers.ParseLocation(lngs[i])
	}

	return list
}

func address(req *http.Request) []string {
	return req.URL.Query()["addr"]
}
