package main

import (
	"io"
)

type apiClient interface {
	SendRequest(path, method string, body io.Reader) ([]byte, error)
	Close()
}
