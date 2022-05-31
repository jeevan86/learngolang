package main

import (
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/ebpf/iovisor/gobpf"
)

func main() {
	gobpf.Start()
	cmd.WaitForSig()
}
