package cilium

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/jeevan86/learngolang/pkg/ebpf/cilium/kprobe"
	"github.com/jeevan86/learngolang/pkg/ebpf/cilium/tracepoint"
	"github.com/jeevan86/learngolang/pkg/log"
	"os"
	"time"
)

var logger = log.NewLogger()

func Start() {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		logger.Fatal("%v", err)
		os.Exit(-1)
	}

	var kpPrograms *kprobe.KpPrograms
	var err error
	kpPrograms, err = kprobe.NewKpSysFuncPrograms()
	if err != nil {
		return
	}
	err = kpPrograms.Start()
	if err != nil {
		return
	}

	var tpPrograms *tracepoint.TPPrograms
	tpPrograms, err = tracepoint.NewTPPrograms()
	if err != nil {
		return
	}
	err = tpPrograms.Start()
	if err != nil {
		return
	}

	go func() {
		defer kpPrograms.Close()
		// Read loop reporting the total amount of times the kernel
		// function was entered, once per second.
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		logger.Info("Waiting for events..")

		for range ticker.C {
			var kpOut = make(map[ebpf.MapType]interface{})
			for _, v := range kpPrograms.GetPrograms() {
				kpPrograms.Lookup(v, kpOut)
				//if err = kpPrograms.Lookup(v, kpOut); err != nil {
				//	logger.Error("reading map: %v", err)
				//} else {
				//	logger.Info("%s called %d times.", v.Name, value)
				//}
			}
			var tpOut = make(map[ebpf.MapType]interface{})
			for _, v := range tpPrograms.GetPrograms() {
				tpPrograms.Lookup(v, tpOut)
				//var value uint64
				//if err = tpPrograms.Lookup(v, &value); err != nil {
				//	logger.Error("reading map: %v", err)
				//} else {
				//	logger.Info("%s called %d times.", v.Name, value)
				//}
			}
		}
	}()
}
