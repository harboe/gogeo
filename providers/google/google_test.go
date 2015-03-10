package google

import (
	"os"
	"io/ioutil"
	"log"
	"testing"

	"github.com/harboe/gogeo/providers"
)

var (
	addressTests = map[string]providers.Result {
		"alekistevej 203, vanlose": providers.Result {
			Address: "Ålekistevej 203, 2720 Vanløse, Denmark",
			Location: providers.Location{
				Latitude: 55.694639,
				Longitude: 12.4796647,
			},
		},
	}
	locationTests = map[string]providers.Result {
		"55.694639,12.4796647": providers.Result {
			Address: "Ålekistevej 203, 2720 Vanløse, Denmark",
			Location: providers.Location{
				Latitude: 55.694639,
				Longitude: 12.4796647,
			},
		},
	}
)

func TestGeoServiceAddress(t *testing.T) {
	g := geoService{}

	for test, expected := range addressTests {
		actual, _ := g.Address(test)

		if actual.Address != expected.Address {
			t.Errorf("expected \"%s\" got \"%s\"", expected.Address, actual.Address)
		}

		if actual.Location != expected.Location {
			t.Errorf("expected \"%v\" got \"%v\"", expected.Location, actual.Location)
		}
	}
}

func TestGeoServiceLocation(t *testing.T) {
	g := geoService{}

	for test, expected := range locationTests {
		actual, _ := g.Location(providers.ParseLocation(test))

		if actual.Address != expected.Address {
			t.Errorf("expected \"%s\" got \"%s\"", expected.Address, actual.Address)
		}

		if actual.Location != expected.Location {
			t.Errorf("expected \"%v\" got \"%v\"", expected.Location, actual.Location)
		}
	}
}

func TestMapServiceAddress(t *testing.T) {
	m := mapService{}
	b, err := m.Address([]string{"alekistevej 203","vigerslev alle 77, valby"}, providers.DefaultMapOptions)

	log.Println("byte:", len(b))
	log.Println("err:", err)

	ioutil.WriteFile("test.png", b, os.ModePerm)
}
