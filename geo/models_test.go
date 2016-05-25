package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidLocation(t *testing.T) {
	assert.Nil(t, Location{}.Valid())
	assert.Nil(t, Location{Latitude: -90.0, Longitude: -180}.Valid())
	assert.Nil(t, Location{Latitude: 90.0, Longitude: 180}.Valid())
	assert.Equal(t, "latitude out of range", Location{Latitude: -90.1}.Valid().Error())
	assert.Equal(t, "latitude out of range", Location{Latitude: 90.1}.Valid().Error())
	assert.Equal(t, "longitude out of range", Location{Longitude: -180.1}.Valid().Error())
	assert.Equal(t, "longitude out of range", Location{Longitude: 180.1}.Valid().Error())
}

func TestParseLatLng(t *testing.T) {
	l, err := NewLocation("")
	assert.Equal(t, Location{}, l)
	assert.Nil(t, err)

	l, err = NewLocation("123")
	assert.Equal(t, Location{}, l)
	assert.Equal(t, "0,0", l.String())
	assert.Equal(t, "bad format", err.Error())

	l, err = NewLocation("123,456,789")
	assert.Equal(t, Location{}, l)
	assert.Equal(t, "0,0", l.String())
	assert.Equal(t, "bad format", err.Error())

	l, err = NewLocation("abc,def")
	assert.Equal(t, Location{}, l)
	assert.Equal(t, "0,0", l.String())
	assert.Equal(t, "parsing latitude: 'abc' invalid syntax", err.Error())

	l, err = NewLocation("10,def")
	assert.Equal(t, Location{Latitude: 10}, l)
	assert.Equal(t, "10,0", l.String())
	assert.Equal(t, "parsing longitude: 'def' invalid syntax", err.Error())

	l, err = NewLocation("10,20")
	assert.Equal(t, "10,20", l.String())
	assert.Equal(t, Location{Latitude: 10, Longitude: 20}, l)
	assert.Nil(t, err)
}

func TestParseSize(t *testing.T) {
	s, err := NewSize("")
	assert.Equal(t, DefaultMapOptions.Size, s)
	assert.Equal(t, "250x250", s.String())
	assert.Nil(t, err)

	s, err = NewSize("123")
	assert.Equal(t, Size{Width: 123, Height: 123}, s)
	assert.Equal(t, "123x123", s.String())
	assert.Nil(t, err)

	s, err = NewSize("abc")
	assert.Equal(t, Size{}, s)
	assert.Equal(t, "0x0", s.String())
	assert.Equal(t, "parsing size: 'abc' invalid syntax", err.Error())

	s, err = NewSize("abcxdef")
	assert.Equal(t, Size{}, s)
	assert.Equal(t, "0x0", s.String())
	assert.Equal(t, "parsing height: 'abc' invalid syntax", err.Error())

	s, err = NewSize("123xdef")
	assert.Equal(t, Size{Width: 123}, s)
	assert.Equal(t, "123x0", s.String())
	assert.Equal(t, "parsing width: 'def' invalid syntax", err.Error())

	s, err = NewSize("123x456")
	assert.Equal(t, Size{Width: 123, Height: 456}, s)
	assert.Equal(t, "123x456", s.String())
	assert.Nil(t, err)
}
