// Package cmdb
// @Title  用于与CMDB对接进行查询
// @Description  用于与CMDB对接进行查询
// @Author  小卒  2022/08/03 10:57
// @Update  小卒  2022/08/03 10:57
package cmdb

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/types"
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

// IsComponent
// @title       是否是CMDB中的组件实例的IP端口
// @description 是否是CMDB中的组件实例的IP端口
// @auth        小卒    2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
// @return      bool                       "是否有查到数据"
func IsComponent(ipPort *types.IpPort) (*types.IpPortMeta, bool) {
	if meta, exists := GetIpPortMeta(ipPort); exists {
		tags := make(types.Tags, 8)
		tags["compId"] = meta.CiId
		tags["compName"] = meta.CiName
		tags["compTypeId"] = meta.CiTypeId
		tags["compTypeName"] = meta.CiTypeName
		return &types.IpPortMeta{
			IpPort: *ipPort,
			Type:   types.IpPortTypeComponent,
			Tags:   tags,
		}, true
	}
	return nil, false
}

// IsHost
// @title       是否是CMDB中的宿主机的IP和端口
// @description 是否是CMDB中的宿主机的IP和端口
// @auth        小卒    2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
// @return      bool                       "是否有查到数据"
func IsHost(ipPort *types.IpPort) (*types.IpPortMeta, bool) {
	if meta, exists := GetIpMeta(ipPort); exists {
		tags := make(types.Tags, 8)
		tags["hostId"] = meta.CiId
		tags["hostName"] = meta.CiName
		return &types.IpPortMeta{
			IpPort: *ipPort,
			Type:   types.IpPortTypeHost,
			Tags:   tags,
		}, true
	}
	return nil, false
}

// GetIpPortMeta TODO 有没有可能IP端口可以描述多个组件?
// @title       GetIpPortMeta
// @description 获取描述IP端口的信息
// @auth        小卒    2022/08/03 10:57
// @param       ip     string         "ip地址"
// @param       port   int32          "端口号"
// @return      meta   *ResIpPortMeta "描述IP端口的信息"
// @return      exists bool           "是否有值"
func GetIpPortMeta(ipPort *types.IpPort) (*ResIpPortMeta, bool) {
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
func GetIpMeta(ipPort *types.IpPort) (*ResIpPortMeta, bool) {
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
