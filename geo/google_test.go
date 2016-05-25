package geo

import (
	"bufio"
	"net/http"
	"strings"
	"testing"
)

var testcases = map[string]string{
	"https://maps.googleapis.com/maps/api/geocode/json?latlng=55.694639,12.4796647": `{
   "results" : [
      {
         "address_components" : [
            {
               "long_name" : "203",
               "short_name" : "203",
               "types" : [ "street_number" ]
            },
            {
               "long_name" : "Ålekistevej",
               "short_name" : "Ålekistevej",
               "types" : [ "route" ]
            },
            {
               "long_name" : "Vanløse",
               "short_name" : "Vanløse",
               "types" : [ "sublocality_level_1", "sublocality", "political" ]
            },
            {
               "long_name" : "København",
               "short_name" : "København",
               "types" : [ "locality", "political" ]
            },
            {
               "long_name" : "København",
               "short_name" : "København",
               "types" : [ "administrative_area_level_2", "political" ]
            },
            {
               "long_name" : "Denmark",
               "short_name" : "DK",
               "types" : [ "country", "political" ]
            },
            {
               "long_name" : "2720",
               "short_name" : "2720",
               "types" : [ "postal_code" ]
            }
         ],
         "formatted_address" : "Ålekistevej 203, 2720 Vanløse, Denmark",
         "geometry" : {
            "location" : {
               "lat" : 55.694639,
               "lng" : 12.4796647
            },
            "location_type" : "ROOFTOP",
            "viewport" : {
               "northeast" : {
                  "lat" : 55.69598798029149,
                  "lng" : 12.4810136802915
               },
               "southwest" : {
                  "lat" : 55.69329001970849,
                  "lng" : 12.4783157197085
               }
            }
         },
         "place_id" : "ChIJB36gXnpRUkYRMVOvdhHXii8",
         "types" : [ "street_address" ]
      }
   ],
   "status" : "OK"
}
`,
}

type mockGoogleFetcher struct{}

func (f *mockGoogleFetcher) Do(req *http.Request) (*http.Response, error) {
	r := bufio.NewReader(strings.NewReader(testcases[req.URL])
	return http.ReadResponse(r, req)
}

var (
	googleMock, _ = New("google", Config{Fetcher: &mockGoogleFetcher{}})
	addressTests  = map[string]Result{
		"copenhagem": Result{
			Address: "Copenhagen, Denmark",
		},
		"alekistevej 203, vanlose": Result{
			Address: "Ålekistevej 203, 2720 Vanløse, Denmark",
			Location: Location{
				Latitude:  55.694639,
				Longitude: 12.4796647,
			},
		},
	}
	locationTests = map[string]Result{
		"55.694639,12.4796647": Result{
			Address: "Ålekistevej 203, 2720 Vanløse, Denmark",
			Location: Location{
				Latitude:  55.694639,
				Longitude: 12.4796647,
			},
		},
	}
)

func TestGeoServiceAddress(t *testing.T) {
	// g := New("google", Config{Fetcher: &mockGoogleFetcher{}})

	for test, expected := range addressTests {
		actual, _ := googleMock.Address(test)

		if actual.Address != expected.Address {
			t.Errorf("expected \"%s\" got \"%s\"", expected.Address, actual.Address)
		}

		if actual.Location != expected.Location {
			t.Errorf("expected \"%v\" got \"%v\"", expected.Location, actual.Location)
		}
	}
}

func TestGeoServiceLocation(t *testing.T) {
	// g :=

	for test, expected := range locationTests {
		l, _ := NewLocation(test)
		actual, _ := googleMock.Location(l)

		if actual.Address != expected.Address {
			t.Errorf("expected \"%s\" got \"%s\"", expected.Address, actual.Address)
		}

		if actual.Location != expected.Location {
			t.Errorf("expected \"%v\" got \"%v\"", expected.Location, actual.Location)
		}
	}
}

// func TestMapServiceAddress(t *testing.T) {
// 	m := mapService{}
// 	b, err := m.Address([]string{"alekistevej 203","vigerslev alle 77, valby"}, providers.DefaultMapOptions)

// 	log.Println("byte:", len(b))
// 	log.Println("err:", err)

// 	ioutil.WriteFile("test.png", b, os.ModePerm)
// }
