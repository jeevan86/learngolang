package http

import (
	"net/http"
	"strings"
)

type client struct {
	internal *http.Client
}

func newClient() *client {
	return &client{
		internal: &http.Client{
			Transport: http.DefaultTransport,
		},
	}
}

func (c *client) post(url, contentType, body string) (*http.Response, error) {
	return c.internal.Post(url, contentType, strings.NewReader(body))
}

func (c *client) get(url string) (*http.Response, error) {
	return c.internal.Get(url)
}
