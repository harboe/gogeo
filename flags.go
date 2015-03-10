package main

import (
	"fmt"
	// "strings"

	"github.com/harboe/gogeo/providers"
)

type (
	AddressList []string
	LocationList []providers.Location
)

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (a *AddressList) String() string {
    return fmt.Sprint(*a)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (a *AddressList) Set(value string) error {
    *a = append(*a, value)
    return nil
}

func (a *AddressList) Type() string {
	return "muuha..."
}

func (l *LocationList) String() string {
	return fmt.Sprint(*l)
}

func (l *LocationList) Set(value string) error {
	// fmt.Println("location list:", value)
	latlng := providers.ParseLocation(value)
	*l = append(*l, latlng)

   return nil
}

func (l *LocationList) Type() string {
	return "buuha.."
}
