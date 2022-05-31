//go:build linux
// +build linux

// Package kprobe
// This program demonstrates attaching an eBPF program to a kernel symbol.
// The eBPF program will be attached to the start of the sys_execve
// kernel function and prints out the number of times it has been called
// every second.
package kprobe

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/jeevan86/learngolang/pkg/ebpf/cilium/base"
	"github.com/jeevan86/learngolang/pkg/log"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
////go:generate ${GOPATH}/bin/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf ./c/kprobe.c -- -I../headers
//go:generate ${PWD}/generate.sh

var logger = log.NewLogger()

type KpPrograms struct {
	ebpfObjs *bpfObjects
	ebpfMaps base.EbpfMapMap
	programs base.ProgramMap
}

var injectedFunMap = base.InjectedFunMap{
	"sys_execve": {
		Name: "sys_execve",
		Maps: []*base.InjectedMapSpec{
			{
				MapK: 0,
				MapT: ebpf.Array,
			},
		},
	},
	"sys_connect": {
		Name: "sys_connect",
		Maps: []*base.InjectedMapSpec{
			{
				MapK: 1,
				MapT: ebpf.Array,
			},
		},
	},
	"sys_accept": {
		Name: "sys_accept",
		Maps: []*base.InjectedMapSpec{
			{
				MapK: 2,
				MapT: ebpf.Array,
			},
		},
	},
	"vfs_read": {
		Name: "vfs_read",
		Maps: []*base.InjectedMapSpec{
			{
				MapK: 3,
				MapT: ebpf.Array,
			},
			{
				MapK: 3,
				MapT: ebpf.PerfEventArray,
			},
		},
	},
}

func NewKpSysFuncPrograms() (*KpPrograms, error) {
	objects := bpfObjects{} // Load pre-compiled programs and maps into the kernel.
	err := loadBpfObjects(&objects, nil)
	if err != nil {
		logger.Fatal("loading objects: %v", err)
		return nil, err
	}
	// bpfPrograms
	programs := base.ReflectGetPrograms(injectedFunMap, (&objects).bpfPrograms, "kprobe_")
	// bpfMaps
	ebpfMaps := base.ReflectGetEbpfMaps((&objects).bpfMaps)
	return &KpPrograms{
		ebpfObjs: &objects,
		ebpfMaps: ebpfMaps,
		programs: programs,
	}, nil
}

func (k *KpPrograms) Start() error {
	for _, e := range k.programs {
		// Open a Kprobe at the entry point of the kernel function and attach the
		// pre-compiled program. Each time the kernel function enters, the program
		// will increment the execution counter by 1. The read loop below polls this
		// map value once per second.
		kp, err := link.Kprobe(e.Name, e.Prog, nil)
		if err != nil {
			logger.Fatal("Opening kprobe failed: %s", err)
			return err
		}
		e.Link = kp
	}
	return nil
}

func (k *KpPrograms) Close() {
	_ = k.ebpfObjs.bpfPrograms.Close()
}

func (k *KpPrograms) GetPrograms() base.ProgramMap {
	return k.programs
}

func (k *KpPrograms) Lookup(p *base.EbpfProgram, out map[ebpf.MapType]interface{}) {
	for _, spec := range p.Maps {
		mp := k.ebpfMaps[spec.MapT]
		key := spec.MapK
		v := base.Lookup(mp, key)
		logger.Info("Kprobe: %s => %v", p.Name, v)
	}
}
