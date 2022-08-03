package cmdb

import (
	"encoding/json"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/log"
	"io"
	"net/http"
)

var logger = log.NewLogger()

type InstanceMeta struct {
	InstanceId   string `json:"instanceId,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
	ClassId      string `json:"classId,omitempty"`
	ClassName    string `json:"className,omitempty"`
}

var resIpPortMetaMap map[string]*ResIpPortMeta

var h *handler

func Start() {
	cfg := &config{
		server: serverConfig{},
		client: clientConfig{},
	}
	h = newHandlerFromConfig(cfg)
}

func Stop() {

}

func GetIpPortMeta(ip string, port int32) (meta *ResIpPortMeta, exists bool) {
	ipPort := key(ip, port)
	meta, exists = resIpPortMetaMap[ipPort]
	if !exists {
		r, e := h.client.post(h.baseUrl, JSON, ipPort)
		defer closeResponse(r)
		if e != nil {
			logger.Error("Failed to post %s, err => %s", h.baseUrl, e.Error())
			return nil, false
		}
		body, _ := readBody(r)
		instMeta := &InstanceMeta{}
		if err := json.Unmarshal([]byte(body), instMeta); err != nil {
			return nil, false
		}
		resIpPortMetaMap[ipPort] = &ResIpPortMeta{
			CompId:       instMeta.InstanceId,
			CompName:     instMeta.InstanceName,
			CompTypeId:   instMeta.ClassId,
			CompTypeName: instMeta.ClassName,
		}
	}
	meta, exists = resIpPortMetaMap[ipPort]
	return
}

func readBody(res *http.Response) (string, error) {
	bodyLen := res.ContentLength
	resBytes := make([]byte, bodyLen)
	_, err := res.Body.Read(resBytes)
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(resBytes), err
}

func closeResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		_ = r.Body.Close()
	}
}

func key(ip string, port int32) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
