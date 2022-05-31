//go:build linux
// +build linux

// This program demonstrates attaching an eBPF program to a kernel tracepoint.
// The eBPF program will be attached to the page allocation tracepoint and
// prints out the number of times it has been reached. The tracepoint fields
// are printed into /sys/kernel/debug/tracing/trace_pipe.
package tracepoint

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/jeevan86/learngolang/pkg/ebpf/cilium/base"
	"github.com/jeevan86/learngolang/pkg/log"
	"sync"
)

var logger = log.NewLogger()

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate ${PWD}/generate.sh

type TPPrograms struct {
	ebpfObjs *bpfObjects
	ebpfMaps base.EbpfMapMap
	programs base.ProgramMap
}

var injectedFunMap = map[string]*base.InjectedFunSpec{
	"mm_page_alloc": {
		Name: "mm_page_alloc",
		Maps: []*base.InjectedMapSpec{
			{
				MapK: 0,
				MapT: ebpf.Array,
			},
		},
	},
}

func NewTPPrograms() (*TPPrograms, error) {
	objects := bpfObjects{} // Load pre-compiled programs and maps into the kernel.
	err := loadBpfObjects(&objects, nil)
	if err != nil {
		logger.Fatal("loading objects: %v", err)
		return nil, err
	}
	// bpfPrograms
	programs := base.ReflectGetPrograms(injectedFunMap, (&objects).bpfPrograms, "tp_")
	// bpfMaps
	ebpfMaps := base.ReflectGetEbpfMaps((&objects).bpfMaps)
	return &TPPrograms{
		ebpfObjs: &objects,
		ebpfMaps: ebpfMaps,
		programs: programs,
	}, nil
}

func (k *TPPrograms) Start() error {
	for _, e := range k.programs {
		// Open a tracepoint and attach the pre-compiled program. Each time
		// the kernel function enters, the program will increment the execution
		// counter by 1. The read loop below polls this map value once per
		// second.
		// The first two arguments are taken from the following pathname:
		// /sys/kernel/debug/tracing/events/kmem/mm_page_alloc
		kp, err := link.Tracepoint("kmem", e.Name, e.Prog, nil)
		if err != nil {
			logger.Fatal("opening tracepoint: %s", err)
			return err
		}
		e.Link = kp
	}
	return nil
}

func (k *TPPrograms) Close() {
	_ = k.ebpfObjs.bpfPrograms.Close()
}

func (k *TPPrograms) GetPrograms() base.ProgramMap {
	return k.programs
}

var lock = &sync.Mutex{}

func (k *TPPrograms) Lookup(p *base.EbpfProgram, out map[ebpf.MapType]interface{}) {
	for _, spec := range p.Maps {
		mp := k.ebpfMaps[spec.MapT]
		key := spec.MapK
		v := base.Lookup(mp, key)
		logger.Info("TracePoint: %s => %v", p.Name, v)
	}
}
