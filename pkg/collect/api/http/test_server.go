package http

import (
	"encoding/json"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api"
	"github.com/jeevan86/learngolang/pkg/server/http"
	"io/ioutil"
	serv "net/http"
	"unsafe"
)

var context = "/collect"

var body200 = "200"

func PrepareTestServer() {
	http.Register(
		context+"/"+uriCtxNetStatics,
		http.POST,
		func(w serv.ResponseWriter, r *serv.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			req := &api.NetStatics{}
			_ = json.Unmarshal(body, req)
			marshaled, _ := json.Marshal(req)
			fmt.Printf("Received: %s\n", marshaled)
			w.WriteHeader(200)
			_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body200)))
		},
	)

	http.Register(
		context+"/"+uriCtxLocalIpLst,
		http.POST,
		func(w serv.ResponseWriter, r *serv.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			req := string(body)
			fmt.Printf("Received: %s\n", req)
			outBytes, _ := json.Marshal(&api.LocalIpLst{
				Data: api.IpList{
					"172.10.231.101", "172.10.231.102", "172.10.231.103",
				},
			})
			w.WriteHeader(200)
			_, _ = w.Write(outBytes)
		},
	)
}
