package base

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/device"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/util/panics"
	"github.com/jeevan86/learngolang/pkg/util/tm"
	"time"
)

var logger = log.NewLogger()

var IsLocalIp func(ip string) bool

func SetCheckLocalIpFunc(f func(ip string) bool) {
	IsLocalIp = f
}

// ProtocolBatchProcessor 将IP包批量转为输出格式的函数
type ProtocolBatchProcessor func(prev, curr, next ProtocolBatch) map[ProtocolClass]interface{}

// OutputStruct 数据输出结构
type OutputStruct struct {
	Values map[IpVersion]map[ProtocolClass]interface{}
	Bucket int64
}

// PacketProcessor 包处理器结构
type PacketProcessor struct {
	ip4PacketCache PacketCache            // ipv4包的缓存
	ip4Processor   ProtocolBatchProcessor // 将ipv4包批量转为输出格式的函数
	ip6PacketCache PacketCache            // ipv6包的缓存
	ip6Processor   ProtocolBatchProcessor // 将ipv6包批量转为输出格式的函数
	localIpAddress map[string]string      // 本地IP地址的缓存
	localIpGetFunc func() []string        // 本地IP获取的函数
	ticker         *time.Ticker           // 一个定时器，1分钟tick一次
	output         chan *OutputStruct     // 输出的数据从这里拿
}

// NewPacketProcessor
// @title       创建包处理器结构指针
// @description 创建包处理器结构指针
// @auth        小卒     2022/08/03 10:57
// @param       ip4     ProtocolBatchProcessor  "将IP包批量转为输出格式"
// @param       ip6     ProtocolBatchProcessor  "将IP包批量转为输出格式"
// @return      r      *PacketProcessor         "包处理器结构指针"
func NewPacketProcessor(ip4, ip6 ProtocolBatchProcessor) *PacketProcessor {
	return &PacketProcessor{
		ip4PacketCache: NewPacketCache(Ipv4),
		ip4Processor:   ip4,
		ip6PacketCache: NewPacketCache(Ipv6),
		ip6Processor:   ip6,
		localIpAddress: make(map[string]string),
		ticker:         time.NewTicker(time.Minute),  // 一分钟处理一次数据
		output:         make(chan *OutputStruct, 32), // TODO: 批量发送，缓冲区32个批次？
	}
}

// IsLocalIp
// @title       判断是否本地IP
// @description 判断是否本地IP
// @auth        小卒  2022/08/03 10:57
// @param       ip   string  "要判断的IP"
// @return      r    bool    "是否本地"
func (p *PacketProcessor) IsLocalIp(ip string) bool {
	_, ok := p.localIpAddress[ip]
	return ok
}

// Out
// @title       PacketProcessor的输出口
// @description PacketProcessor的输出口
// @auth        小卒   2022/08/03 10:57
// @return      chan  *OutputStruct  "数据库输出的Channel"
func (p *PacketProcessor) Out() chan *OutputStruct {
	return p.output
}

// Process
// @title       IP包处理函数
// @description IP包处理函数，先放入cache
// @auth        小卒     2022/08/03 10:57
// @param       packet  gopacket.Packet  "go-packet包"
func (p *PacketProcessor) Process(packet gopacket.Packet) {
	// 只处理IP包
	if IsIpPacket(packet) {
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
	} else {
		// Check for errors
		if err := packet.ErrorLayer(); err != nil {
			logger.Warn("Error decoding some part of the packet: %s", err.Error())
		}
	}
}

// packetTime
// @title       获得时间戳、时间窗口
// @description 获得时间戳、时间窗口
// @auth        小卒     2022/08/03 10:57
// @param       packet  gopacket.Packet  "go-packet包"
// @return      bucket  int64            "时间窗"
// @return      millis  int64            "时间戳"
// @return      p       gopacket.Packet  "包"
func packetTime(packet gopacket.Packet) (bucket, millis int64, p gopacket.Packet) {
	millis = time.Now().UnixMilli()
	bucket = tm.TruncToMinuteTs(millis / 1000)
	p = packet
	return
}

// Start
// @title       启动PacketProcessor包处理器结构
// @description 启动PacketProcessor包处理器结构
// @auth        小卒     2022/08/03 10:57
// @param       localIpGetFunc func() []string  "获取本地ip的函数"
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

// routine
// @title       PacketProcessor周期运行的函数
// @description 启动PacketProcessor后将定时执行
// @auth        小卒  2022/08/03 10:57
// @param       t    time.Time  "time.Ticker的时间"
func (p *PacketProcessor) routine(t time.Time) {
	_, _ = panics.SafeRun(func() {
		p.updateAddresses()
		p.flush(tm.TruncToMinuteTs(t.Unix()))
	})
}

// updateAddresses
// @title       更新本地地址
// @description routine中将调用
// @auth        小卒  2022/08/03 10:57
func (p *PacketProcessor) updateAddresses() {
	// TODO：是否保留一个周期
	devIpList := p.getLocalAddresses()
	internal := make(map[string]string, len(devIpList)*2)
	for _, devIp := range devIpList {
		internal[devIp] = devIp
	}
	if p.localIpGetFunc != nil {
		othIpList := p.localIpGetFunc()
		for _, othIp := range othIpList {
			internal[othIp] = othIp
		}
	}
	p.localIpAddress = internal
}

// getLocalAddresses
// @title       获得本地地址
// @description 获得本地地址
// @auth        小卒  2022/08/03 10:57
// @return      r1   []string   "ip列表"
func (p *PacketProcessor) getLocalAddresses() []string {
	// 更新本地的IP地址
	addresses := make([]string, 0)
	for _, inf := range device.AllDevices() {
		for _, address := range inf.Addresses {
			addresses = append(addresses, address.IP.String())
		}
	}
	return addresses
}

// flush
// @title       flush周期运行的函数
// @description flush周期运行的函数
// @auth        小卒  2022/08/03 10:57
// @param       m    int64   "时间窗口的时间戳"
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

// Stop
// @title       停止PacketProcessor
// @description 停止PacketProcessor的ticker
// @auth        小卒  2022/08/03 10:57
func (p *PacketProcessor) Stop() {
	p.ticker.Stop()
}
