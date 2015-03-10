package providers

import (
	"errors"
	"strings"
)

var (
	geoProviders = map[string]GeoServiceConstructor{}
)

func Geo(name, apikey string) (GeoService, error) {
	if c, ok := geoProviders[name]; ok {
		return c(apikey)
	}

	return nil, errors.New("no provided found named: " + name +
		". valid providers are: " +
		strings.Join(GeoProviders(), ","))
}

func RegisterGeo(name string, constuctor GeoServiceConstructor) {
	geoProviders[name] = constuctor
	// fmt.Println("geo providers:", geoProviders, "is", geoProviders[name])
}

func GeoProviders() []string {
	list := []string{}

	for k, _ := range geoProviders {
		list = append(list, k)
	}

	return list
}
