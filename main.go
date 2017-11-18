package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"fmt"

	"os"

	"encoding/json"

	"github.com/czerwonk/ovirt_api"
)

const version string = "0.1.0"

var (
	showVersion   = flag.Bool("version", false, "Prints version info")
	listenAddress = flag.String("listen-address", ":11337", "Address to listen for web service requests")
	user          = flag.String("username", "user@internal", "API username")
	pass          = flag.String("password", "", "API password")
	apiUrl        = flag.String("api-url", "https://ovirt.engine/ovirt-engine/api", "API url")
	insecure      = flag.Bool("insecure", false, "Skip SSL verification")
	templateFile  = flag.String("template", "ovirt_vm_template.xml", "Template file path")
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

	http.HandleFunc("/", errorHandler(handleRequest))
	log.Println("ovirt-zero-touch " + version)
	log.Println("Server started. Start listening on " + *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func printVersion() {
	fmt.Println("ovirt-zero-touch")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			log.Fatalln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	vm := Request{}
	err = json.Unmarshal(bytes, &vm)
	if err != nil {
		return err
	}

	b, err := createVm(&vm)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func createVm(vm *Request) ([]byte, error) {
	client, err := ovirt_api.NewClient(*apiUrl, *user, *pass, *insecure, &logger{})
	if err != nil {
		return nil, err
	}

	body, err := getVmCreateRequest(vm)
	if err != nil {
		return nil, err
	}

	b, err := client.SendRequest("vms", "POST", body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func getVmCreateRequest(vm *Request) (io.Reader, error) {
	w := &bytes.Buffer{}

	b, err := ioutil.ReadFile(*templateFile)
	if err != nil {
		return w, err
	}

	tmpl, err := template.New("create-vm").Parse(string(b))
	if err != nil {
		return w, err
	}

	tmpl.Execute(w, vm)
	return w, nil
}