package actuator

import (
	"encoding/json"
	"github.com/jeevan86/learngolang/pkg/server/http"
	"io/ioutil"
	serv "net/http"
	"unsafe"
)

var loggers = context + "/loggers"

type LogLevel struct {
	Prefix string `json:"prefix"`
	Level  string `json:"level"`
}

var body200 = "200"

func Init() {
	http.Register(
		loggers+"/update",
		http.POST,
		func(w serv.ResponseWriter, r *serv.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			logLevel := new(LogLevel)
			_ = json.Unmarshal(body, logLevel)
			logger.SetLevels(logLevel.Prefix, logLevel.Level)
			w.WriteHeader(200)
			_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body200)))
		},
	)

	http.Register(
		loggers+"/select",
		http.POST,
		func(w serv.ResponseWriter, r *serv.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			logLevel := new(LogLevel)
			_ = json.Unmarshal(body, logLevel)
			outBytes, _ := json.Marshal(logger.GetLevels(logLevel.Prefix))
			w.WriteHeader(200)
			_, _ = w.Write(outBytes)
		},
	)
}
