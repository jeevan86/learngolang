package http

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/config"
	"net/http"
)

var server = &http.Server{}

func Start() {
	address := config.GetConfig().Server.Address
	go func() {
		for _, handler := range handlers {
			http.Handle(handler.path, handler)
		}
		server.Addr = address
		server.Handler = nil
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}

func Stop() {
	_ = server.Close()
}
