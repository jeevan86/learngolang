package main

import (
	"github.com/jeevan86/learngolang/cmd/util"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc"
)

func main() {
	grpc.PrepareTestServer()
	util.WaitForSig()
}
