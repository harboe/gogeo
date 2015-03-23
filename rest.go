package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/harboe/gogeo/providers"
)

const helpTemplate = `
<html>
<head>
	<title>gogeo help</title>
</head>
<body>

</body>
</html>`

func RestService(port string) {
	router := httprouter.New()
	router.GET("/v1/:name/png", imgHandler)
	router.GET("/v1/:name/json", geoHandler)
	router.GET("/v1/:name/xml", geoHandler)
	router.GET("/v1/:name/txt", geoHandler)
	// router.NotFound = helpHandler

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
	log.Println("here...")
	geo, err := providers.Geo(ps.ByName("name"))

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
		log.Println("address:", a)
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
	geo, err := providers.Geo(ps.ByName("name"))

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

	for _, p := range providers.Providers() {
		help := `
			` + p + `
		`

		w.Write([]byte(help))
	}

	w.Write([]byte("documentation not implemented"))
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
