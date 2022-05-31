package main

import (
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc"
)

func main() {
	grpc.PrepareTestServer()
	cmd.WaitForSig()
}
