package capture

import (
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/util/parallel"
	"github.com/jeevan86/learngolang/pkg/util/reactor"
)

// captureHandeFunc
// @title       包装一下包处理函数
// @description 包装一下包处理函数，提供不同的并发模型
// @auth        小卒     2022/08/03 10:57
// @param       process func(packet gopacket.Packet)  "包处理函数"
// @return      r1      func(uint64, gopacket.Packet) "包装后的包处理函数"
func captureHandeFunc(process func(packet gopacket.Packet)) func(uint64, gopacket.Packet) {
	switch conf.ParType {
	case config.ParTypeRoutine:
		return parRoutinesFunc(process) // with Routines
	case config.ParTypeReactor:
		return reactorFluxFunc(process) // with RxGo
	default:
		return parRoutinesFunc(process) // with Routines
	}
}

// parRoutinesFunc
// @title       多goroutine并发模型
// @description 包装一下包处理函数，提供多goroutine并发模型
// @auth        小卒     2022/08/03 10:57
// @param       process func(packet gopacket.Packet)  "包处理函数"
// @return      r1      func(uint64, gopacket.Packet) "包装后的包处理函数"
func parRoutinesFunc(process func(packet gopacket.Packet)) func(uint64, gopacket.Packet) {
	parallelism := conf.Routine.Parallelism
	chBuffSz := conf.Routine.ChBufferSize
	isShareCh := conf.Routine.ShareChan
	wrappedFn := func(p interface{}) {
		process(p.(gopacket.Packet))
	}
	// with Routines
	routines := parallel.NewParRoutines(
		parallelism,
		chBuffSz,
		isShareCh,
		wrappedFn,
	)
	return func(i uint64, p gopacket.Packet) {
		routines.Dispatch(i, p)
	}
}

// reactorFluxFunc
// @title       reactor响应式并发模型
// @description 包装一下包处理函数，提供reactor响应式并发模型（依赖rxgo而实现）
// @auth        小卒     2022/08/03 10:57
// @param       process func(packet gopacket.Packet)  "包处理函数"
// @return      r1      func(uint64, gopacket.Packet) "包装后的包处理函数"
func reactorFluxFunc(process func(packet gopacket.Packet)) func(uint64, gopacket.Packet) {
	bufferSz := conf.Reactor.BufferSz
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
