package main

import (
	"github.com/jeevan86/learngolang/cmd/util"
	"github.com/jeevan86/learngolang/pkg/collect/api/http"
	server "github.com/jeevan86/learngolang/pkg/server/http"
)

func main() {
	http.PrepareTestServer()
	server.Start()
	util.WaitForSig()
	server.Stop()
}
