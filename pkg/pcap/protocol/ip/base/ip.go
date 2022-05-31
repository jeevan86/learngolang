package base

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/device"
	"github.com/jeevan86/learngolang/pkg/util/panics"
	"github.com/jeevan86/learngolang/pkg/util/tm"
	"time"
)

var logger = log.NewLogger()

var IsLocalIp func(ip string) bool

func SetCheckLocalIpFunc(f func(ip string) bool) {
	IsLocalIp = f
}

type IpVersion int8

const (
	NotIp IpVersion = -1
	Ipv4  IpVersion = 1
	Ipv6  IpVersion = 2
)

type ProtocolClass int8

const (
	ProtocolUnknown ProtocolClass = -1
	ProtocolIgmp    ProtocolClass = 1
	ProtocolIcmp    ProtocolClass = 2
	ProtocolTcp     ProtocolClass = 3
	ProtocolUdp     ProtocolClass = 4
)

func IsIpPacket(packet gopacket.Packet) bool {
	return Version(packet) > 0
}

func Version(packet gopacket.Packet) IpVersion {
	if nil != packet.Layer(layers.LayerTypeIPv4) {
		return Ipv4
	}
	if nil != packet.Layer(layers.LayerTypeIPv6) {
		return Ipv6
	}
	return NotIp
}

type LayerIp interface {
	GetSrcIp() string
	GetDstIp() string
	GetPktSz() uint16
}

type ProtocolBatchProcessor func(prev, curr, next ProtocolBatch) map[ProtocolClass]interface{}

type OutputStruct struct {
	Values map[IpVersion]map[ProtocolClass]interface{}
	Bucket int64
}

type PacketProcessor struct {
	ip4PacketCache PacketCache
	ip4Processor   ProtocolBatchProcessor
	ip6PacketCache PacketCache
	ip6Processor   ProtocolBatchProcessor
	localIpAddress map[string]string
	localIpGetFunc func() []string
	ticker         *time.Ticker
	output         chan *OutputStruct
}

func NewPacketProcessor(ip4, ip6 ProtocolBatchProcessor) *PacketProcessor {
	return &PacketProcessor{
		ip4PacketCache: NewPacketCache(Ipv4),
		ip4Processor:   ip4,
		ip6PacketCache: NewPacketCache(Ipv6),
		ip6Processor:   ip6,
		localIpAddress: make(map[string]string),
		ticker:         time.NewTicker(time.Minute), // 一分钟处理一次数据
		output:         make(chan *OutputStruct, 32),
	}
}

func (p *PacketProcessor) IsLocalIp(ip string) bool {
	_, ok := p.localIpAddress[ip]
	return ok
}

func (p *PacketProcessor) Out() chan *OutputStruct {
	return p.output
}

// Process 处理IP包
func (p *PacketProcessor) Process(packet gopacket.Packet) {
	v := Version(packet)
	if v == NotIp {
		logger.Warn("Not an ip packet.")
		return
	} else if v == Ipv4 {
		ip4 := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		if ip4.SrcIP.IsLoopback() || ip4.DstIP.IsLoopback() {
			return
		}
		p.ip4PacketCache.PutPacket(packetTime(packet))
	} else if v == Ipv6 {
		ip6 := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv6)
		if ip6.SrcIP.IsLoopback() || ip6.DstIP.IsLoopback() {
			return
		}
		p.ip6PacketCache.PutPacket(packetTime(packet))
	}
}

func packetTime(packet gopacket.Packet) (bucket, millis int64, p gopacket.Packet) {
	millis = time.Now().UnixMilli()
	bucket = tm.TruncToMinuteTs(millis / 1000)
	p = packet
	return
}

func (p *PacketProcessor) Start(localIpGetFunc func() []string) {
	p.localIpGetFunc = localIpGetFunc
	go func() {
		for {
			t, ok := <-p.ticker.C
			if !ok {
				break
			}
			p.routine(t)
		}
	}()
}

func (p *PacketProcessor) routine(t time.Time) {
	panics.SafeRun(func() {
		p.updateAddresses()
		p.flush(tm.TruncToMinuteTs(t.Unix()))
	})
}

func (p *PacketProcessor) updateAddresses() {
	// TODO：是否保留一个周期
	devIpList := getLocalAddresses()
	internal := make(map[string]string, len(devIpList)*2)
	for _, ip := range devIpList {
		internal[ip] = ip
	}
	if p.localIpGetFunc != nil {
		othIpList := p.localIpGetFunc()
		for _, ip := range othIpList {
			internal[ip] = ip
		}
	}
	p.localIpAddress = internal
}

func getLocalAddresses() []string {
	// 更新本地的IP地址
	addresses := make([]string, 0)
	for _, inf := range device.AllDevices() {
		for _, address := range inf.Addresses {
			addresses = append(addresses, address.IP.String())
		}
	}
	return addresses
}

func (p *PacketProcessor) flush(m int64) {
	// 处理上一个1分钟时间片内的数据
	keyNext := m
	keyCurr := keyNext - 60
	keyPrev := keyCurr - 60
	out := make(map[IpVersion]map[ProtocolClass]interface{}, 2)
	out[Ipv4] = p.flushIp4(keyNext, keyCurr, keyPrev)
	out[Ipv6] = p.flushIp6(keyNext, keyCurr, keyPrev)
	p.output <- &OutputStruct{
		Values: out,
		Bucket: m,
	}
}

func (p *PacketProcessor) Stop() {
	p.ticker.Stop()
}
