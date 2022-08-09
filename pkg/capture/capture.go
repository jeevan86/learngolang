package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/util/id"
	"github.com/jeevan86/learngolang/pkg/util/parallel"
	"github.com/jeevan86/learngolang/pkg/util/reactor"
	"os"
	"time"
)

var logger = log.NewLogger()

var conf = config.GetConfig().Agent.Capture

// StartCapture
// @title       启动pcap采集器
// @description 启动一个基于go-packet和pcap库的网卡流量采集器
// @auth        小卒     2022/08/03 10:57
// @param       process func(packet gopacket.Packet)     "包处理器"
func StartCapture(process func(packet gopacket.Packet)) {
	for _, device := range conf.Devices {
		handle, err := createHandle(device.Duration, device.Prefix, device.Snaplen, device.Promisc)
		if err != nil {
			// permission issue
			logger.Fatal("Unable to create handler, cause: %s", err.Error())
			os.Exit(-1)
		} else {
			go func() { startCapture(handle, process) }()
		}
	}
	logger.Info("Packet capture started.")
}

// createHandle
// @title       创建一个设备流量处理器
// @description 创建一个设备流量处理器
// @auth        小卒     2022/08/03 10:57
// @param       duration string       "时长(最小精度毫秒)"
// @param       device   string       "网卡设备"
// @param       snaplen  string       "snaplen?"
// @param       promisc  bool         "promisc?"
// @return      r1       *pcap.Handle "go-packet的包接收器"
// @return      r2       error        "错误信息"
func createHandle(duration, device string, snaplen int32, promisc bool) (*pcap.Handle, error) {
	// 设置1秒？
	secs, _ := time.ParseDuration(duration)
	//handle, err := pcap.OpenLive("en0", 40, true, secs)
	h, err := pcap.OpenLive(device, snaplen, promisc, secs)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// startCapture
// @title       启动pcap采集器
// @description 从包处理器得到数据，并分配数据的id，然后交给函数处理
// @auth        小卒     2022/08/03 10:57
// @param       handle *pcap.Handle                   "包接收器"
// @param       process func(packet gopacket.Packet)  "包处理函数"
func startCapture(handle *pcap.Handle, process func(packet gopacket.Packet)) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	_, idx, after := id.CyclingIdxFunc()
	captured := captureHandeFunc(process)
	for packet := range packetSource.Packets() {
		captured(*idx, packet)
		after()
	}
}

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
