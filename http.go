package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/harboe/gogeo/geo"
)

func RestService(port string) {
	router := httprouter.New()
	router.GET("/:name/png", imgHandler)
	router.GET("/:name/json", geoHandler)
	router.GET("/:name/xml", geoHandler)
	router.GET("/:name/yml", geoHandler)

	fmt.Println("route=GET /:name/:format[png,json,xml,yml]")

	routeHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		router.ServeHTTP(w, req)
	})

	chain := alice.New(loggingHandler, corsHandlers).Then(routeHandler)

	log.Fatal(http.ListenAndServe(port, chain))
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func corsHandlers(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Stop here if its Preflighted OPTIONS request
		if req.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func geoHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	qry := req.URL.Query()
	provider, err := config.NewWithKey(ps.ByName("name"), qry.Get("key"))
	// provider, err := geo.New(ps.ByName("name"), geo.Config{APIKey: qry.Get("key")})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var results []geo.Result

	for _, l := range location(qry) {
		if r, err := provider.Location(l); err == nil {
			results = append(results, r)
		}
	}

	for _, a := range address(qry) {
		log.Println("address:", a)
		if r, err := provider.Address(a); err == nil {
			results = append(results, r)
		}
	}

	_, pretty := qry["pretty"]
	format := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	b, err := Marshal(format, &results, pretty)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func imgHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	qry := req.URL.Query()
	geo, err := config.NewWithKey(ps.ByName("name"), qry.Get("key"))

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	opts := mapOptions(qry)
	markers := address(qry)

	for _, loc := range location(qry) {
		markers = append(markers, loc.String())
	}

	b, err := geo.Image(markers, opts)

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func mapOptions(qry url.Values) (opts geo.MapOptions) {
	opts.Size, _ = geo.NewSize(qry.Get("size"))

	if z := qry.Get("zoom"); len(z) > 0 {
		opts.Zoom, _ = strconv.ParseUint(z, 0, 10)
	}

	if s := qry.Get("scale"); len(s) > 0 {
		opts.Scale, _ = strconv.ParseUint(s, 0, 10)
	}

	return
}

func location(qry url.Values) (res []geo.Location) {
	lngs := qry["loc"]

	for i := 0; i < len(lngs); i++ {
		if loc, err := geo.NewLocation(lngs[i]); err == nil {
			res = append(res, loc)
		}
	}

	return
}

func address(qry url.Values) []string {
	return qry["addr"]
}
