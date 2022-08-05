package cmdb

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"github.com/reactivex/rxgo/v2"
	"time"
)

// resIpPortMetaCache ristretto是意大利超浓咖啡
var resIpPortMetaCache *ristretto.Cache // map[string]*ResIpPortMeta
const resIpPortMetaItemCost = 4
const resIpPortMetaItemLive = 30 * time.Minute

const rxChBufferTime = 5 * time.Minute
const rxChBufferCount = 32

var rxCh = make(chan rxgo.Item, 256)

func initTtlCache() {
	// TODO: 过期的元素，到时候一起刷新
	onEvict := func(item *ristretto.Item) { rxCh <- rxgo.Of(item.Value) }
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		OnEvict:     onEvict,
	})
	if err != nil {
		panic(err)
	}
	resIpPortMetaCache = cache
}

func startEvictedRefresher() {
	// TODO: 缓存过期数据的刷新
	observable := rxgo.FromChannel(rxCh)
	observable.
		Map(
			func(ctx context.Context, i interface{}) (interface{}, error) {
				obj := i.(*ResIpPortMeta).IpPort
				return &obj, nil
			},
		).
		BufferWithTimeOrCount(
			rxgo.WithDuration(rxChBufferTime),
			rxChBufferCount,
		).
		DoOnNext(
			func(buf interface{}) {
				items := buf.([]*backend.IpPort)
				for k, v := range *client.pollInstanceMetaList(items) {
					if v != nil {
						transAndCache(&k, v)
					}
				}
			},
		)
}

func stopEvictedRefresher() {
	close(rxCh)
}

func transAndCache(k *backend.IpPort, instMeta *InstanceMeta) *ResIpPortMeta {
	// 转换一下
	meta := toResIpPortMeta(instMeta)
	// 缓存一下
	resIpPortMetaCache.SetWithTTL(
		*k, meta, resIpPortMetaItemCost, resIpPortMetaItemLive)
	return meta
}
