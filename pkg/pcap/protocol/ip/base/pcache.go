package base

import (
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/pkg/util/panics"
)

type DistinctPacketId struct {
	DstIp string
	PktId uint32
}

type PacketItem struct {
	Id     DistinctPacketId
	Millis int64
	Packet gopacket.Packet
}

type PacketBatch map[DistinctPacketId]*PacketItem
type ProtocolBatch map[ProtocolClass]PacketBatch

type PacketCache interface {
	PutPacket(int64, int64, gopacket.Packet)
	GetBatch(int64) ProtocolBatch
	DelBatch(int64)
	distinct(gopacket.Packet) DistinctPacketId
	protocol(gopacket.Packet) ProtocolClass
}

func NewPacketCache(version IpVersion) PacketCache {
	return cacheCreator[version](newPacketCache())
}

var cacheCreator = make(map[IpVersion]func(*DefaultPacketCache) PacketCache)

func newPacketCache() *DefaultPacketCache {
	c := &DefaultPacketCache{
		cache: make(map[int64]ProtocolBatch),
		ch:    make(chan *CacheCommand, 2048),
	}
	c.startSyncRoutine()
	return c
}

type cacheCommand uint8

const (
	cacheCmdDel cacheCommand = 0
	cacheCmdAdd cacheCommand = 1
)

type CacheCommand struct {
	cmd      cacheCommand
	bucket   int64
	protocol ProtocolClass
	item     *PacketItem
}

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

func (c *DefaultPacketCache) startSyncRoutine() {
	go func() {
		for {
			cmd, ok := <-c.ch
			if !ok {
				break
			}
			panics.SafeRun(func() { c.syncRoutine(cmd) })
		}
	}()
}

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
