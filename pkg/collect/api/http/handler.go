package http

import (
	"encoding/json"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api"
	"github.com/jeevan86/learngolang/pkg/util/str"
	"io"
	"net/http"
)

const (
	uriCtxNetStatics = "netstatics"
	uriCtxLocalIpLst = "localip"
)

const (
	JSON = "application/json"
	TEXT = "text/plain"
)

type Handler struct {
	baseUrl string
	client  *client
}

func NewHandler(serverAddr string) *Handler {
	return &Handler{
		baseUrl: serverAddr,
		client:  newClient(),
	}
}

var emptyString = str.EMPTY()

func (h *Handler) Save(msg *api.NetStatics) (*string, error) {
	bytes, _ := json.Marshal(msg)
	body := string(bytes)
	url := h.baseUrl + "/" + uriCtxNetStatics
	res, err := h.client.post(url, JSON, body)
	defer closeResponse(res)
	if err != nil {
		return &emptyString, err
	}
	bodyLen := res.ContentLength
	resBytes := make([]byte, bodyLen)
	_, err = res.Body.Read(resBytes)
	if err != nil && err != io.EOF {
		return &emptyString, err
	}
	result := string(resBytes)
	return &result, nil
}

func (h *Handler) LocalIp(nodeIp string) (*string, error) {
	url := h.baseUrl + "/" + uriCtxLocalIpLst
	res, err := h.client.post(url, JSON, fmt.Sprintf("{\"nodeIp\": \"%s\"}", nodeIp))
	defer closeResponse(res)
	if err != nil {
		return &emptyString, err
	}
	bodyLen := res.ContentLength
	resBytes := make([]byte, bodyLen)
	_, err = res.Body.Read(resBytes)
	if err != nil && err != io.EOF {
		return &emptyString, err
	}
	result := string(resBytes)
	return &result, nil
}

func closeResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		_ = r.Body.Close()
	}
}
