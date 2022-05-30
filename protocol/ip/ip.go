package ip

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	logging "gopackettest/logger"
)

var logger = logging.LoggerFactory.NewLogger([]string{"stdout"}, []string{"stderr"})

var localIpAddresses *map[string]pcap.InterfaceAddress

func GetAddressInfo(ip string) pcap.InterfaceAddress {
	return (*localIpAddresses)[ip]
}

func UpdateAddresses(addresses []pcap.InterfaceAddress) {
	internal := make(map[string]pcap.InterfaceAddress, len(addresses)*2)
	for _, address := range addresses {
		internal[address.IP.String()] = address
	}
	localIpAddresses = &internal
}

// Metric 指标的结构
type Metric struct {
	Metric   string `json:"metric"`
	Instance string `json:"instance_name"`
	//"protocol": "tcp"
	//"protocol_ver": ""
	Tags      map[string]interface{} `json:"metric_tags"`
	Timestamp uint64                 `json:"timestamp"`
	Value     float64                `json:"value"`
}

const (
	NotIp = -1
	Ipv4  = 1
	Ipv6  = 2
)

func IsIpPacket(packet gopacket.Packet) bool {
	return version(packet) > 0
}

func version(packet gopacket.Packet) int {
	if nil != packet.Layer(layers.LayerTypeIPv4) {
		return Ipv4
	}
	if nil != packet.Layer(layers.LayerTypeIPv6) {
		return Ipv6
	}
	return NotIp
}

func ProcessIpPacket(packet gopacket.Packet) {
	v := version(packet)
	if v == NotIp {
		logger.Warn("Not an ip packet.")
		return
	} else if v == Ipv4 {
		processIp4Packet(packet)
	} else if v == Ipv6 {
		processIp6Packet(packet)
	}
}
