package providers

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/xml"
)

var DefaultMapOptions = MapOptions{Size{250, 250}, 0, 0}

type (
	Location struct {
		Latitude  float64 `json:"lat" xml:"location>lat"`
		Longitude float64 `json:"lng" xml:"location>lng"`
	}
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
	Results []Result
	Size    struct {
		Width  uint64
		Height uint64
	}
	MapOptions struct {
		Size
		Scale uint64
		Zoom  uint64
	}
	GeoService interface {
		Location(loc Location) (Result, error)
		Address(address string) (Result, error)
		Static(markers []string, options MapOptions) ([]byte, error)
	}
	GeoServiceConstructor func(key string) (GeoService, error)
)

func (s Size) String() string {
	return fmt.Sprintf("%vx%v", s.Width, s.Height)
}

func (l Location) String() string {
	return fmt.Sprintf("%v,%v", l.Latitude, l.Longitude)
}

func ParseSize(size string) Size {
	if len(size) > 0 {
		if strings.Index(size, "x") == -1 {
			if size, err := strconv.ParseUint(size, 0, 10); err == nil {
				return Size{size, size}
			}
		} else {
			arr := strings.Split(size, "x")
			h, _ := strconv.ParseUint(arr[0], 0, 10)
			w, _ := strconv.ParseUint(arr[1], 0, 10)
			return Size{w, h}
		}
	}

	return Size{250, 250}
}

func ParseLocation(loc string) Location {
	// fmt.Println("loc:", loc)

	if len(loc) == 0 {
		return Location{}
	}

	arr := strings.Split(loc, ",")
	lat, _ := strconv.ParseFloat(arr[0], 10)
	lng, _ := strconv.ParseFloat(arr[1], 10)

	return Location{lat, lng}
}

func (r Results) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "results",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, value := range r {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: "result"},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}
