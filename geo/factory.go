package geo

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	googleGeoURL   = "https://maps.googleapis.com/maps/api/geocode/json"
	googleImgURL   = "https://maps.googleapis.com/maps/api/staticmap"
	bingGeoURL     = "https://dev.virtualearth.net/REST/v1/Locations"
	bingImgURL     = ""
	mapquestGeoURL = "https://open.mapquestapi.com/geocoding/v1/"
	mapquestImgURL = ""
)

// Config represents optional provider configurations
type Config struct {
	// Fetcher to use when getting geo request
	// If nothing is specificed it will default back to http.DefaultClient
	Fetcher Fetcher
	// APIKey to use in the geo service. If nothing is specificed it will
	// try to find an Env variables name GOGEO_{name}
	APIKey string
}

// Providers return a list of available providers
func Providers() []string {
	return []string{"bing", "google", "mapquest"}
}

// New returns new instance of the provider specificed by name
func New(name string, opts ...Config) (Provider, error) {
	var cfg Config

	if len(opts) == 0 {
		cfg = Config{}
	} else {
		cfg = opts[0]
	}

	if cfg.Fetcher == nil {
		cfg.Fetcher = http.DefaultClient
	}

	if len(cfg.APIKey) == 0 {
		cfg.APIKey = APIKey(name)
	}

	switch strings.ToLower(name) {
	case "google":
		return &googleAPI{Config: cfg, Geo: googleGeoURL, Img: googleImgURL}, nil
	case "bing":
		return &bingAPI{Config: cfg, Geo: bingGeoURL, Img: bingImgURL}, nil
	case "mapquest":
		return &mapquestAPI{Config: cfg, Geo: mapquestGeoURL, Img: mapquestImgURL}, nil
	}

	return nil, fmt.Errorf("not found: %s", name)
}

func APIKey(key string) string {
	return os.Getenv(fmt.Sprintf("GOGEO_%s", strings.ToUpper(key)))
}
