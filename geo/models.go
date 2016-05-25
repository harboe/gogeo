package geo

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	// DefaultSize if nothing is specificed in the request
	DefaultSize = Size{Width: 250, Height: 250}
	// DefaultMapOptions if nothing is specificed in the request
	DefaultMapOptions = MapOptions{Size: DefaultSize, Scale: 0, Zoom: 0}
)

type (
	// Location represents geo latitude,longitude coordinates
	Location struct {
		Latitude  float64 `json:"lat" xml:"location>lat"`
		Longitude float64 `json:"lng" xml:"location>lng"`
	}
	// Result represents the response for any given provider
	Result struct {
		Query    string `json:"query" xml:"query,attr"`
		Address  string `json:"address" xml:"address"`
		Street   string `json:"street,omitempty" xml:"street,omitempty"`
		Country  string `json:"country" xml:"country"`
		City     string `json:"city,omitempty" xml:"city,omitempty"`
		Zip      string `json:"zip,omitempty" xml:"zip,omitempty"`
		State    string `json:"state,omitempty" xml:"state,omitempty"`
		Location `json:"location"`
	}
	// Size of a image
	Size struct {
		Width  uint64
		Height uint64
	}
	// MapOptions need to create static images.
	MapOptions struct {
		Size
		Scale uint64
		Zoom  uint64
	}
	// Provider for geo and reveresed address lookup
	Provider interface {
		Location(loc Location) (Result, error)
		Address(address string) (Result, error)
		Image(markers []string, options MapOptions) ([]byte, error)
	}
)

func (s Size) String() string {
	return fmt.Sprintf("%vx%v", s.Width, s.Height)
}

// Valid returns a error if location is out of range.
func (l Location) Valid() error {
	if !(l.Latitude >= -90.0 && l.Latitude <= 90.0) {
		return fmt.Errorf("latitude out of range")
	}

	if !(l.Longitude >= -180.0 && l.Longitude <= 180.0) {
		return fmt.Errorf("longitude out of range")
	}

	return nil
}

func (l Location) String() string {
	return fmt.Sprintf("%v,%v", l.Latitude, l.Longitude)
}

// NewSize converts eighter {width}x{height} or {size} to a Size.
func NewSize(size string) (s Size, err error) {
	size = strings.TrimSpace(size)

	if len(size) == 0 {
		return DefaultSize, nil
	}

	// only a single value is specificed
	if strings.Index(size, "x") == -1 {
		if s.Height, err = strconv.ParseUint(size, 0, 10); err != nil {
			return s, fmt.Errorf("parsing size: '%s' invalid syntax", size)
		}
		s.Width = s.Height
	} else {
		arr := strings.Split(size, "x")
		if s.Width, err = strconv.ParseUint(arr[0], 0, 10); err != nil {
			return s, fmt.Errorf("parsing height: '%s' invalid syntax", arr[0])
		}
		if s.Height, err = strconv.ParseUint(arr[1], 0, 10); err != nil {
			return s, fmt.Errorf("parsing width: '%s' invalid syntax", arr[1])
		}
	}

	return
}

// NewLocation converts eighter {lat},{lng} to Location.
func NewLocation(loc string) (l Location, err error) {
	if len(loc) == 0 {
		return l, nil
	}

	arr := strings.Split(loc, ",")

	if len(arr) != 2 {
		return l, fmt.Errorf("bad format")
	}

	if l.Latitude, err = strconv.ParseFloat(arr[0], 10); err != nil {
		return l, fmt.Errorf("parsing latitude: '%s' invalid syntax", arr[0])
	}
	if l.Longitude, err = strconv.ParseFloat(arr[1], 10); err != nil {
		return l, fmt.Errorf("parsing longitude: '%s' invalid syntax", arr[1])
	}

	return l, l.Valid()
}
