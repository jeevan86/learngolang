package cmdb

import (
	"encoding/json"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"io"
	"net/http"
	"strings"
)

const (
	JSON = "application/json"
	TEXT = "text/plain"
)

type httpClient struct {
	internal *http.Client
}

func (c *httpClient) post(url, contentType, body string) (*http.Response, error) {
	return c.internal.Post(url, contentType, strings.NewReader(body))
}

func (c *httpClient) get(url string) (*http.Response, error) {
	return c.internal.Get(url)
}

// readBody
// @title       readBody
// @description 获取http响应的body字符串
// @auth        小卒    2022/08/03 10:57
// @param       res    *http.Response "ip地址"
// @return      ""     string         "响应体字符串"
// @return      ""     error          "错误"
func readBody(res *http.Response) (string, error) {
	bodyLen := res.ContentLength
	resBytes := make([]byte, bodyLen)
	_, err := res.Body.Read(resBytes)
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(resBytes), err
}

// closeResponse
func closeResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		_ = r.Body.Close()
	}
}

type cmdbClient struct {
	baseUrl    string
	httpClient *httpClient
}

// pollInstanceMeta 请求CMDB的接口，获取信息。TODO: 需要接口文档
// @title       pollInstanceMeta
// @description 获取描述IP端口的CMDB对象的信息
// @auth        小卒    2022/08/03 10:57
// @param       ip     string         "ip地址"
// @param       port   int32          "端口号"
// @return      meta   *InstanceMeta  "CMDB对象"
func (c *cmdbClient) pollInstanceMeta(ipPort *backend.IpPort) *InstanceMeta {
	key := *ipPort
	r, e := c.httpClient.post(c.baseUrl, JSON, fmt.Sprintf("%s:%d", key.Ip, key.Port))
	defer closeResponse(r)
	if e != nil {
		logger.Error("Failed to post %s, err => %s", c.baseUrl, e.Error())
		return nil
	}
	body, _ := readBody(r)
	instMeta := &InstanceMeta{}
	err := json.Unmarshal([]byte(body), instMeta)
	if err != nil {
		return nil
	}
	return instMeta
}

// pollInstanceMeta 请求CMDB的接口，获取信息。TODO: 需要接口文档
// @title       pollInstanceMeta
// @description 获取描述IP端口的CMDB对象的信息
// @auth        小卒    2022/08/03 10:57
// @param       ip     string         "ip地址"
// @param       port   int32          "端口号"
// @return      meta   *InstanceMeta  "CMDB对象"
func (c *cmdbClient) pollInstanceMetaList(ipPort []*backend.IpPort) *map[backend.IpPort]*InstanceMeta {
	js, _ := json.Marshal(ipPort)
	reqBody := string(js)
	r, e := c.httpClient.post(c.baseUrl, JSON, reqBody)
	defer closeResponse(r)
	if e != nil {
		logger.Error("Failed to post %s, err => %s", c.baseUrl, e.Error())
		return nil
	}
	body, _ := readBody(r)
	instMetaMap := new(map[backend.IpPort]*InstanceMeta)
	err := json.Unmarshal([]byte(body), instMetaMap)
	if err != nil {
		return nil
	}
	return instMetaMap
}

var client *cmdbClient

func newClient(serverAddr string) *cmdbClient {
	return &cmdbClient{
		baseUrl: serverAddr,
		httpClient: &httpClient{
			internal: &http.Client{
				Transport: http.DefaultTransport,
			},
		},
	}
}

func newClientFromConfig(cfg *config) {
	client = newClient(cfg.server.url)
}
