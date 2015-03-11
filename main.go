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
	format       string
	addrList     AddressList
	locList      LocationList
	size         string
	zoom         uint64
	scale        uint64
	pretty       bool
	apikey       string
	data         string
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

	for _, provider := range providers.GeoProviders() {
		c := &cobra.Command{
			Use:     provider,
			Short:   provider + " provider",
			Long:    provider + " provider",
			Example: "$ gogeo " + provider + " -a \"vigerslev alle 77, valby\" -l 55.694639,12.4796647",
			Run:     runProvider,
		}
		imgCmd := &cobra.Command{
			Use:     "img",
			Short:   "Generator static map",
			Long:    "Generator static map",
			Example: "$ gogeo " + provider + " img -a \"vigerslev alle 77, valby\" test.png",
			Run:     runImageProvider,
		}
		c.AddCommand(imgCmd)
		f := c.Flags()
		p := c.PersistentFlags()

		p.VarP(
			&addrList, "address", "a", "addresses")
		p.VarP(
			&locList, "location", "l", "latitude,longitude")
		f.StringVarP(
			&data, "data", "d", "", "data input file")
		f.StringVar(
			&format, "format", "json", "output format json|xml|txt")
		f.BoolVar(
			&pretty, "pretty", false, "pretty print")
		p.StringVar(
			&apikey, "key", "", "optional depending on the specific provider")
		imgCmd.Flags().StringVar(
			&size, "size", "250x250", "map size use for png")
		imgCmd.Flags().Uint64Var(
			&scale, "scale", 0, "usage")
		imgCmd.Flags().Uint64Var(
			&zoom, "zoom", 0, "map zoom level, varies depending provider")

		rootCmd.AddCommand(c)
		providerKeys[provider] = serverCmd.Flags().String(provider+"-key", "", "optional depending on the specific provider")
	}
	rootCmd.Execute()
}

func restService(cmd *cobra.Command, args []string) {
	fmt.Printf("rest service ready at http://%s\n", port)
	RestService(port)
}

func runImageProvider(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		fmt.Println("error: need output filename")
		return
	}

	geo, err := providers.Geo(cmd.Parent().Use, apikey)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	opts := providers.MapOptions{
		Size:  providers.ParseSize(size),
		Zoom:  zoom,
		Scale: scale,
	}
	markers := addrList

	for _, loc := range locList {
		markers = append(markers, loc.String())
	}

	if b, err := geo.Static(markers, opts); err != nil {
		fmt.Println(err)
	} else {
		ioutil.WriteFile(args[0]+".png", b, os.ModePerm)
	}
}

func runProvider(cmd *cobra.Command, args []string) {
	geo, err := providers.Geo(cmd.Use, apikey)
	v := providers.Results{}

	if err != nil {
		fmt.Println(err.Error())
		return
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
		if len(args) > 0 {
			ioutil.WriteFile(args[0]+"."+format, b, os.ModePerm)
		} else {
			fmt.Println(string(b))
		}
	}
}
