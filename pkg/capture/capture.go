package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/util/id"
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
