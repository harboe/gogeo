package main

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"
)

type (
	Encoding interface {
		Marshal(v interface{}, pretty bool) ([]byte, error)
	}
	jsonEncoding struct{}
	xmlEncoding struct{}
	txtEncoding struct{}
)

var (
	encoders = map[string]Encoding{
		"json": jsonEncoding{},
		"xml": xmlEncoding{},
		"txt": txtEncoding{},
	}
)

func (e jsonEncoding) Marshal(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(&v, "", "  ")
	} else {
		return json.Marshal(&v)
	}
}

func (e xmlEncoding) Marshal(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return xml.MarshalIndent(&v, "", "  ")
	} else {
		return xml.Marshal(&v)
	}
}

func (e txtEncoding) Marshal(v interface{}, pretty bool) ([]byte, error) {
	return yaml.Marshal(&v)
}

func Marshal(ct string, v interface{}, pretty bool) ([]byte, error) {
	e, ok := encoders[ct];

	if !ok {
		e = encoders["json"]
	}

	return e.Marshal(v, pretty)
}