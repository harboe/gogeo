package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/harboe/gogeo/providers"
	_ "github.com/harboe/gogeo/providers/bing"
	_ "github.com/harboe/gogeo/providers/google"
)

var (
	port         string
	file         string
	format       string
	addrList     AddressList
	locList      LocationList
	size         string
	zoom         uint64
	scale        uint64
	pretty       bool
	apikey       string
	providerKeys = map[string]*string{}
)

func main() {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "run as a server",
		Long:  "gogeo: as a rest service",
		Run:   restService,
	}
	serverCmd.Flags().StringVarP(&port, "port", "p", "localhost:8080", "server listing port")

	rootCmd := &cobra.Command{
		Use:  "gogeo",
		Long: "muuha...",
	}
	rootCmd.AddCommand(serverCmd)

	for _, p := range providers.GeoProviders() {
		c := &cobra.Command{
			Use:     p,
			Short:   p + " provider",
			Long:    p + " provider",
			Example: "$ gogeo " + p + " -a \"vigerslev alle 77, valby\" -l 55.694639,12.4796647",
			Run:     provider,
		}
		f := c.Flags()

		f.VarP(
			&addrList, "address", "a", "addresses")
		f.VarP(
			&locList, "location", "l", "latitude,longitude")
		f.StringVarP(
			&file, "file", "f", "", "output file")
		f.StringVar(
			&format, "format", "json", "output format json|xml|txt|png")
		f.BoolVar(
			&pretty, "pretty", false, "pretty print")
		f.StringVar(
			&apikey, "key", "", "optional depending on the specific provider")
		f.StringVar(
			&size, "size", "250x250", "map size use for png")
		f.Uint64Var(
			&scale, "scale", 1, "usage")
		f.Uint64Var(
			&zoom, "zoom", 1, "map zoom level, varies depending provider")

		rootCmd.AddCommand(c)
		providerKeys[p] = serverCmd.Flags().String(p+"-key", "", "optional depending on the specific provider")
	}
	rootCmd.Execute()
}

func restService(cmd *cobra.Command, args []string) {
	fmt.Printf("rest service ready at http://%s\n", port)
	RestService(port)
}

func provider(cmd *cobra.Command, args []string) {
	geo, err := providers.Geo(cmd.Use, apikey)
	v := providers.Results{}

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, addr := range addrList {
		if a, err := geo.Address(addr); err == nil {
			v = append(v, a)
		} else {
			fmt.Println("skip:", addr, "error:", err.Error())
		}
	}

	for _, loc := range locList {
		if l, err := geo.Location(loc); err == nil {
			v = append(v, l)
		} else {
			fmt.Println("skip:", loc, "error:", err.Error())
		}
	}

	if b, err := Marshal(format, &v, pretty); err != nil {
		fmt.Println("marshal error:", err)
	} else {
		if len(file) > 0 {
			ioutil.WriteFile(file, b, os.ModePerm)
		} else {
			fmt.Println(string(b))
		}
	}
}
