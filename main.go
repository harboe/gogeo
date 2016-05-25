package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/harboe/gogeo/geo"
	"github.com/spf13/cobra"
)

var (
	server struct {
		Port string
	}
	image    imageFlags
	format   formatFlags
	addrList addressList
	locList  locationList
	config   configFlags
)

func main() {
	serverCmd := &cobra.Command{
		Use:   "http",
		Short: "execute a httpserver",
		Long:  "gogeo: as a rest service",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rest service ready at http://%s\n", server.Port)
			RestService(server.Port)
		},
	}
	serverCmd.Flags().StringVarP(&server.Port, "port", "p", "localhost:8080", "server listing port")

	envCmd := &cobra.Command{
		Use:   "env",
		Short: "display env apikeys",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range geo.Providers() {
				fmt.Printf("GOGEO_%s = %s\n",
					strings.ToUpper(p), geo.APIKey(p))
			}
		},
	}

	rootCmd := &cobra.Command{
		Use:  "gogeo",
		Long: "Awesome geo fetching and backend service",
	}
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "")
	rootCmd.AddCommand(serverCmd, envCmd)

	for _, provider := range geo.Providers() {
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

		f.BoolVarP(&format.Json, "json", "j", false, "output json format")
		f.BoolVarP(&format.Yaml, "yml", "y", false, "output yml format")
		f.BoolVarP(&format.Xml, "xml", "x", false, "output xml format")
		f.BoolVarP(&format.Pretty, "pretty", "p", false, "pretty print")

		imgCmd.Flags().StringVar(
			&image.Size, "size", "250x250", "map size use for png")
		imgCmd.Flags().Uint64Var(
			&image.Scale, "scale", 0, "usage")
		imgCmd.Flags().Uint64Var(
			&image.Zoom, "zoom", 0, "map zoom level, varies depending provider")

		rootCmd.AddCommand(c)

		serverCmd.Flags().String(provider+"-key", "", "optional depending on the specific provider")

		p.StringVar(&config.APIKey, "key", "", "optional depending on the specific provider")
	}
	rootCmd.Execute()
}

func runImageProvider(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		fmt.Println("error: need output filename")
		return
	}

	provider, err := config.New(cmd.Parent().Use)

	if err != nil {
		fmt.Printf("provider: %v\n", err)
		return
	}

	opts, err := image.Map()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	markers := addrList

	for _, loc := range locList {
		markers = append(markers, loc.String())
	}

	if b, err := provider.Image(markers, opts); err != nil {
		fmt.Println(err)
	} else {

		ioutil.WriteFile(args[0]+".png", b, os.ModePerm)
	}
}

func runProvider(cmd *cobra.Command, args []string) {
	if len(addrList) == 0 && len(locList) == 0 {
		cmd.Help()
		return
	}

	provider, err := config.New(cmd.Use)
	v := []geo.Result{}

	if err != nil {
		fmt.Printf("provider: %v\n", err)
		return
	}

	for _, addr := range addrList {
		if a, err := provider.Address(addr); err == nil {
			v = append(v, a)
		} else {
			fmt.Println("skip:", addr, "error:", err.Error())
		}
	}

	for _, loc := range locList {
		if l, err := provider.Location(loc); err == nil {
			v = append(v, l)
		} else {
			fmt.Println("skip:", loc, "error:", err.Error())
		}
	}

	b, err := format.Marshal(&v)

	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	if len(args) > 0 {
		ioutil.WriteFile(format.Filename(args[0]), b, os.ModePerm)
	} else {
		fmt.Println(string(b))
	}
}
