package cmdb

import (
	"net/http"
	"strings"
)

const (
	JSON = "application/json"
	TEXT = "text/plain"
)

type handler struct {
	baseUrl string
	client  *client
}

func newHandlerFromConfig(cfg *config) *handler {
	return newHandler(cfg.server.url)
}

func newHandler(serverAddr string) *handler {
	return &handler{
		baseUrl: serverAddr,
		client:  newClient(),
	}
}

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
