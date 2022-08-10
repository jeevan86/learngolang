package base

import (
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/pkg/util/panics"
)

// DistinctPacketId 根据IP和包ID可以确定是同一个包
type DistinctPacketId struct {
	DstIp string
	PktId uint32
}

// PacketItem 包含包唯一ID、时间戳毫秒、具体的包
type PacketItem struct {
	Id     DistinctPacketId
	Millis int64
	Packet gopacket.Packet
}

// PacketBatch 用于根据DistinctPacketId进行排重
type PacketBatch map[DistinctPacketId]*PacketItem

// ProtocolBatch 根据协议进行合并
type ProtocolBatch map[ProtocolClass]PacketBatch

// PacketCache 包缓存，用于排重、按时间窗口聚合等
type PacketCache interface {
	PutPacket(int64, int64, gopacket.Packet)
	GetBatch(int64) ProtocolBatch
	DelBatch(int64)
	distinct(gopacket.Packet) DistinctPacketId
	protocol(gopacket.Packet) ProtocolClass
}

// cacheCreator 根据IpVersion进行分类的缓存构建器
var cacheCreator = make(map[IpVersion]func(*DefaultPacketCache) PacketCache)

// NewPacketCache
// @title       NewPacketCache
// @description 使用cacheCreator创建缓存
// @auth        小卒     2022/08/03 10:57
// @param       version IpVersion   "ip6或者ip4"
// @return      r       PacketCache "缓存实现"
func NewPacketCache(version IpVersion) PacketCache {
	return cacheCreator[version](newPacketCache())
}

// newPacketCache
// @title       newPacketCache
// @description 创建一个默认的缓存结构
// @auth        小卒     2022/08/03 10:57
// @return      r       *DefaultPacketCache "缓存实现"
func newPacketCache() *DefaultPacketCache {
	c := &DefaultPacketCache{
		cache: make(map[int64]ProtocolBatch),
		ch:    make(chan *CacheCommand, 2048),
	}
	c.startSyncRoutine()
	return c
}

// cacheCommand 操作缓存的指令id，用于同步操作
type cacheCommand uint8

const (
	cacheCmdDel cacheCommand = 0 // 删除指令
	cacheCmdAdd cacheCommand = 1 // 增加指令
)

// CacheCommand 操作缓存的指令，用于同步操作
type CacheCommand struct {
	cmd      cacheCommand  // 指令id
	bucket   int64         // 桶编号
	protocol ProtocolClass // ip子协议
	item     *PacketItem   // ip包信息
}

// DefaultPacketCache 默认的缓存结构，实现了一些通用的内部操作
type DefaultPacketCache struct {
	cache map[int64]ProtocolBatch
	ch    chan *CacheCommand
}

func (c *DefaultPacketCache) putPacket(bucket, millis int64, id DistinctPacketId, pro ProtocolClass, p gopacket.Packet) {
	c.ch <- &CacheCommand{
		cmd:      cacheCmdAdd,
		bucket:   bucket,
		protocol: pro,
		item: &PacketItem{
			Id:     id,
			Millis: millis,
			Packet: p,
		},
	}
}

func (c *DefaultPacketCache) getBatch(bucket int64) ProtocolBatch {
	return c.cache[bucket]
}

func (c *DefaultPacketCache) delBatch(bucket int64) {
	c.ch <- &CacheCommand{
		cmd:    cacheCmdDel,
		bucket: bucket,
		item:   nil,
	}
}

func (c *DefaultPacketCache) newProtocolBatch(bucket int64) ProtocolBatch {
	c.cache[bucket] = make(ProtocolBatch, 6)
	//"golang.org/x/tools/go/types/typeutil"
	//typeutil.MakeHasher()
	return c.cache[bucket]
}

func (c *DefaultPacketCache) newPacketBatch(bucket int64, clz ProtocolClass) PacketBatch {
	c.cache[bucket][clz] = make(PacketBatch, 10240)
	return c.cache[bucket][clz]
}

// startSyncRoutine
// @title       startSyncRoutine
// @description 启动缓存的运作
// @auth        小卒     2022/08/03 10:57
func (c *DefaultPacketCache) startSyncRoutine() {
	go func() {
		for {
			cmd, ok := <-c.ch
			if !ok {
				break
			}
			_, _ = panics.SafeRun(func() { c.syncRoutine(cmd) })
		}
	}()
}

// syncRoutine
// @title       syncRoutine
// @description 同步地执行内部缓存指令
// @auth        小卒  2022/08/03 10:57
// @param       cmd  *CacheCommand   "缓存操作指令"
func (c *DefaultPacketCache) syncRoutine(cmd *CacheCommand) {
	if cmd.cmd == cacheCmdDel {
		delete(c.cache, cmd.bucket)
	} else if cmd.cmd == cacheCmdAdd {
		bucket := cmd.bucket
		batch := c.getBatch(bucket)
		if batch == nil {
			batch = c.newProtocolBatch(bucket)
		}
		items := batch[cmd.protocol]
		if items == nil {
			items = c.newPacketBatch(bucket, cmd.protocol)
		}
		item := cmd.item
		items[item.Id] = item
	}
}
