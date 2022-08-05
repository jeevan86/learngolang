// Package cmdb
// @Title  用于与CMDB对接进行查询
// @Description  用于与CMDB对接进行查询
// @Author  小卒  2022/08/03 10:57
// @Update  小卒  2022/08/03 10:57
package cmdb

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"github.com/jeevan86/learngolang/pkg/log"
)

var logger = log.NewLogger()

// Start
// @title       初始化并启动包
// @description 创建客户端、缓存等
// @auth        小卒 2022/08/03 10:57
func Start() {
	// TODO: CMDB服务器客户端的配置
	cfg := &config{
		server: serverConfig{},
		client: clientConfig{},
	}
	newClientFromConfig(cfg)
	initTtlCache()
	startEvictedRefresher()
}

func Stop() {
	stopEvictedRefresher()
}

// GetIpPortMeta TODO 有没有可能IP端口可以描述多个组件?
// @title       GetIpPortMeta
// @description 获取描述IP端口的信息
// @auth        小卒    2022/08/03 10:57
// @param       ip     string         "ip地址"
// @param       port   int32          "端口号"
// @return      meta   *ResIpPortMeta "描述IP端口的信息"
// @return      exists bool           "是否有值"
func GetIpPortMeta(ipPort *backend.IpPort) (*ResIpPortMeta, bool) {
	var meta *ResIpPortMeta
	cached, exists := resIpPortMetaCache.Get(*ipPort)
	if !exists || cached == nil {
		if instMeta := client.pollInstanceMeta(ipPort); instMeta == nil {
			return nil, false
		} else {
			meta = transAndCache(ipPort, instMeta)
		}
	} else {
		meta = cached.(*ResIpPortMeta)
	}
	return meta, true
}

// GetIpMeta TODO 有没有可能IP端口可以描述多个组件?
// @title       GetIpMeta
// @description 获取描述IP的信息
// @auth        小卒    2022/08/03 10:57
// @param       ip     string         "ip地址"
// @return      meta   *ResIpPortMeta "描述IP端口的信息"
// @return      exists bool           "是否有值"
func GetIpMeta(ipPort *backend.IpPort) (*ResIpPortMeta, bool) {
	var meta *ResIpPortMeta
	cached, exists := resIpPortMetaCache.Get(*ipPort)
	if !exists || cached == nil {
		if instMeta := client.pollInstanceMeta(ipPort); instMeta == nil {
			return nil, false
		} else {
			meta = transAndCache(ipPort, instMeta)
		}
	} else {
		meta = cached.(*ResIpPortMeta)
	}
	return meta, true
}
