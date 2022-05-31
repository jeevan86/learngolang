//go:build linux
// +build linux

package main

import (
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/ebpf/cilium"
)

func main() {
	cilium.Start()
	cmd.WaitForSig()
}
