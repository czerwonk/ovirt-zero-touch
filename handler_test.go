package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mock struct {
	t      *testing.T
	called bool
	body   string
	closed bool
	err    error
}

func (m *mock) SendRequest(path, method string, body io.Reader) ([]byte, error) {
	m.called = true

	if path != "vms" {
		m.t.Fatalf("expected call with path 'vms', got %s", path)
	}

	if method != "POST" {
		m.t.Fatalf("expected call with method POST, got %s", method)
	}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		m.t.Fatal("could not read body")
	}
	m.body = strings.TrimSpace(string(b))

	return []byte{}, m.err
}

func (m *mock) Close() {
	m.closed = true
}

func TestHandlerCreatesVM(t *testing.T) {
	m := runTest(t, nil, http.StatusOK)

	if !m.called {
		t.Fatal("no request was sent")
	}

	expected := expectedBody()
	if m.body != expected {
		t.Fatalf("got unexpected body\nexpected:\n%s\n\ngot:\n%s", expected, m.body)
	}
}

func expectedBody() string {
	b, err := ioutil.ReadFile("tests/body.xml")
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(b))
}

func TestHandlerClosesClient(t *testing.T) {
	m := runTest(t, nil, http.StatusOK)

	if !m.closed {
		t.Fatal("handler does not close api connection")
	}
}

func runTest(t *testing.T, err error, expectedStatus int) *mock {
	m := &mock{t: t, err: err}
	api := func() (apiClient, error) {
		return m, nil
	}
	h := newHandler(api, loadTestTemplate)
	srv := httptest.NewServer(h)
	defer srv.Close()

	b := loadTestRequest()
	res, err := http.Post(srv.URL, "", b)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if res.StatusCode != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, res.Status)
	}

	return m
}

func loadTestRequest() io.Reader {
	b, err := ioutil.ReadFile("examples/request.json")
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(b)
}

func loadTestTemplate() ([]byte, error) {
	return ioutil.ReadFile("examples/template.xml")
}
