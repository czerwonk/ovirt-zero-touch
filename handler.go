package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/template"
)

type apiClientFunc func() (apiClient, error)
type templateFunc func() ([]byte, error)

type handler struct {
	apiFunc      apiClientFunc
	templateFunc templateFunc
}

func newHandler(api apiClientFunc, template templateFunc) http.Handler {
	h := &handler{api, template}

	r := http.NewServeMux()
	r.HandleFunc("/", errorHandler(h.handleRequest))
	return r
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *handler) handleRequest(w http.ResponseWriter, r *http.Request) error {
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
	vm.Memory *= 1048576 // MB -> Bytes

	b, err := h.createVM(&vm)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func (h *handler) createVM(vm *Request) ([]byte, error) {
	client, err := h.apiFunc()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	body, err := h.getVMCreateRequest(vm)
	if err != nil {
		return nil, err
	}

	b, err := client.SendRequest("vms", "POST", body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h *handler) getVMCreateRequest(vm *Request) (io.Reader, error) {
	b, err := h.templateFunc()
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("create-vm").Parse(string(b))
	if err != nil {
		return nil, err
	}

	w := &bytes.Buffer{}
	err = tmpl.Execute(w, vm)
	return w, err
}
