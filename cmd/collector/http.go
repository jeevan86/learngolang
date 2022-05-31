package main

import (
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/collect/api/http"
	server "github.com/jeevan86/learngolang/pkg/server/http"
)

func main() {
	http.PrepareTestServer()
	server.Start()
	cmd.WaitForSig()
	server.Stop()
}
