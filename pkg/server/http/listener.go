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
		for _, h := range handlers {
			http.Handle(h.path, h)
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
