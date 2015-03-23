package providers

import (
	"errors"
	"os"
	"strings"
)

var (
	geoProviders    = map[string]GeoServiceConstructor{}
	geoProviderKeys = map[string][]*string{}
)

func Geo(name string) (GeoService, error) {
	if c, ok := geoProviders[name]; ok {
		return c(getApiKey(name))
	}

	return nil, errors.New("no provided found named: " + name +
		". valid providers are: " +
		strings.Join(Providers(), ","))
}

func Register(name string, constuctor GeoServiceConstructor) {
	geoProviders[name] = constuctor
}

func Providers() []string {
	list := []string{}

	for k, _ := range geoProviders {
		list = append(list, k)
	}

	return list
}

func AddApiKey(name string, keyP *string) {
	geoProviderKeys[name] = append(geoProviderKeys[name], keyP)
}

func getApiKey(name string) string {
	for _, p := range geoProviderKeys[name] {
		if p != nil && len(*p) > 0 {
			return *p
		}
	}

	return os.Getenv(strings.ToUpper("gogeo_" + name + "_key"))
}
