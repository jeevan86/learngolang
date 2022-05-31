package main

import (
	"fmt"
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/collect/api/http"
	server "github.com/jeevan86/learngolang/pkg/server/http"
	"testing"
)

func Test_Main(t *testing.T) {
	http.PrepareTestServer()
	server.Start()
	collector.Start(packetProcessor.Out())
	packetProcessor.Start(localIpGetFunc)
	startCapture(parRoutinesFunc()) // with Routines
	// startCapture(reactorFluxFunc()) // with RxGo
	fmt.Println("packet capture started.")
	cmd.WaitForSig()
	packetProcessor.Stop()
	collector.Stop()
	server.Stop()
}
