//go:build linux
// +build linux

package main

import (
	"github.com/jeevan86/learngolang/cmd/util"
	"github.com/jeevan86/learngolang/pkg/ebpf/cilium"
)

func main() {
	cilium.Start()
	util.WaitForSig()
}
