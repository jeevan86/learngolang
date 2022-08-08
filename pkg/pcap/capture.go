package pcap

import (
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/util/parallel"
	"github.com/jeevan86/learngolang/pkg/util/reactor"
)

func captureHandeFunc(process func(packet gopacket.Packet)) func(uint64, gopacket.Packet) {
	switch config.GetConfig().Agent.Pcap.ParType {
	case config.ParTypeRoutine:
		return parRoutinesFunc(process) // with Routines
	case config.ParTypeReactor:
		return reactorFluxFunc(process) // with RxGo
	default:
		return parRoutinesFunc(process) // with Routines
	}
}

func parRoutinesFunc(process func(packet gopacket.Packet)) func(uint64, gopacket.Packet) {
	parallelism := config.GetConfig().Agent.Pcap.Routine.Parallelism
	// with Routines
	routines := parallel.NewParRoutines(
		parallelism,
		2048,
		false,
		func(p interface{}) {
			process(p.(gopacket.Packet))
		},
	)
	return func(i uint64, p gopacket.Packet) {
		routines.Dispatch(i, p)
	}
}

func reactorFluxFunc(process func(packet gopacket.Packet)) func(uint64, gopacket.Packet) {
	bufferSz := config.GetConfig().Agent.Pcap.Reactor.BufferSz
	fluxSink := reactor.NewFluxSink(bufferSz)
	fluxSink.Map(func(e interface{}) interface{} {
		return e
		//}).FlatMap(func(e interface{}) *util.FluxSink {
		//	packet := e.(gopacket.Packet)
		//	newFluxSink := util.NewFluxSink(2048)
		//	newFluxSink.Next(packet)
		//	return newFluxSink
	}).DoOnNext(func(p interface{}) {
		process(p.(gopacket.Packet))
	}).Subscribe(func(e interface{}) {
		// do nothing
	})
	return func(i uint64, p gopacket.Packet) {
		fluxSink.Next(p)
	}
}
