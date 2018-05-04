package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/czerwonk/ovirt_api/api"
)

const version string = "0.2.3"

var (
	showVersion   = flag.Bool("version", false, "Prints version info")
	listenAddress = flag.String("listen-address", ":11337", "Address to listen for web service requests")
	user          = flag.String("username", "user@internal", "API username")
	pass          = flag.String("password", "", "API password")
	apiURL        = flag.String("api-url", "https://ovirt.engine/ovirt-engine/api", "API url")
	insecure      = flag.Bool("insecure", false, "Skip SSL verification")
	templateFile  = flag.String("template", "ovirt_vm_template.xml", "Template file path")
	debug         = flag.Bool("debug", false, "Enables verbose output")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: ovirt-zero-touch [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	log.Println("ovirt-zero-touch " + version)

	h := newHandler(newAPIClient, loadTemaplate)

	log.Println("Server started. Start listening on " + *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, h))
}

func printVersion() {
	fmt.Println("ovirt-zero-touch")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
}

func loadTemaplate() ([]byte, error) {
	return ioutil.ReadFile(*templateFile)
}

type apiClientAdapter struct {
	*api.Client
}

func newAPIClient() (apiClient, error) {
	opts := []api.ClientOption{}
	if *insecure {
		opts = append(opts, api.WithInsecure())
	}

	if *debug {
		opts = append(opts, api.WithDebug())
	}

	c, err := api.NewClient(*apiURL, *user, *pass, opts...)
	if err != nil {
		return nil, err
	}

	return &apiClientAdapter{c}, nil
}
