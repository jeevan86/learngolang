package actuator

import (
	"encoding/json"
	"fmt"
	logging "gopackettest/logger"
	"io/ioutil"
	"net/http"
	"unsafe"
)

var logger = logging.LoggerFactory.NewLogger([]string{"stdout"}, []string{"stderr"})

type LogLevel struct {
	Prefix string `json:"prefix"`
	Level  string `json:"level"`
}

var body404 = "404"
var body200 = "200"

type actuatorLoggerUpdate string

func (h actuatorLoggerUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		logLevel := new(LogLevel)
		_ = json.Unmarshal(body, logLevel)
		logger.SetLevels(logLevel.Prefix, logLevel.Level)
		w.WriteHeader(200)
		_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body200)))
	} else {
		w.WriteHeader(404)
		_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body404)))
	}
}

type actuatorLoggerSelect string

func (h actuatorLoggerSelect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		logLevel := new(LogLevel)
		_ = json.Unmarshal(body, logLevel)
		outBytes, _ := json.Marshal(logger.GetLevels(logLevel.Prefix))
		w.WriteHeader(200)
		_, _ = w.Write(outBytes)
	} else {
		w.WriteHeader(404)
		_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body404)))
	}
}

func StartActuator() {
	go func() {
		http.Handle("/actuator/loggers/update", actuatorLoggerUpdate("update"))
		http.Handle("/actuator/loggers/select", actuatorLoggerSelect("select"))
		err := http.ListenAndServe("127.0.0.1:8630", nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}
