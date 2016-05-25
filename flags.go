package main

import (
	"fmt"
	"net/http"
	// "strings"

	"github.com/harboe/gogeo/geo"
	"github.com/harboe/gogeo/geo/middleware"
)

type (
	addressList  []string
	locationList []geo.Location

	imageFlags struct {
		Size  string
		Zoom  uint64
		Scale uint64
	}
	formatFlags struct {
		Yaml   bool
		Json   bool
		Xml    bool
		Pretty bool
	}
	configFlags struct {
		Verbose bool
		APIKey  string
	}
)

func (c configFlags) New(name string) (geo.Provider, error) {
	return c.NewWithKey(name, c.APIKey)
}

func (c configFlags) NewWithKey(name, key string) (geo.Provider, error) {
	var f geo.Fetcher

	if c.Verbose {
		f = middleware.Logger(http.DefaultClient)
	}

	return geo.New(name, geo.Config{APIKey: key, Fetcher: f})

}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (a *addressList) String() string {
	return fmt.Sprint(*a)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (a *addressList) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (a *addressList) Type() string {
	return "address list type!!?" // no idear what i should return
}

func (l *locationList) String() string {
	return fmt.Sprint(*l)
}

func (l *locationList) Set(value string) error {
	// fmt.Println("location list:", value)
	latlng, err := geo.NewLocation(value)
	*l = append(*l, latlng)

	return err
}

func (l *locationList) Type() string {
	return "location list type!!?" // no idear what i should return
}

func (i imageFlags) Map() (opt geo.MapOptions, err error) {
	if opt.Size, err = geo.NewSize(i.Size); err != nil {
		return
	}
	opt.Zoom = i.Zoom
	opt.Scale = i.Scale
	return
}

func (f formatFlags) String() string {
	switch {
	case f.Json:
		return "json"
	case f.Yaml:
		return "yml"
	case f.Xml:
		return "xml"
	}

	return "json"
}

func (f formatFlags) Marshal(v interface{}) ([]byte, error) {
	return Marshal(f.String(), v, f.Pretty)
}

func (f formatFlags) Filename(name string) string {

	return ""
}
