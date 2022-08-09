package main

import (
	"github.com/jeevan86/learngolang/cmd/util"
	"github.com/jeevan86/learngolang/pkg/ebpf/iovisor/gobpf"
)

func main() {
	gobpf.Start()
	util.WaitForSig()
}
