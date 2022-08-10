package device

import (
	"github.com/google/gopacket/pcap"
	"github.com/jeevan86/learngolang/pkg/log"
)

var logger = log.NewLogger()

// AllDevices
// @title       本地所有的网卡设备
// @description 本地所有的网卡设备
// @auth        小卒  2022/08/03 10:57
// @return      r    []string   "网卡设备列表"
func AllDevices() []pcap.Interface {
	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		logger.Fatal(err.Error())
	}
	return devices
}
