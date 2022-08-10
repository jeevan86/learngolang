package cmdb

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/types"
)

// InstanceMeta 接口传输对象-组件基本属性
type InstanceMeta struct {
	Ip           string `json:"ip,omitempty"`
	Port         int32  `json:"port,omitempty"`
	InstanceId   string `json:"instanceId,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
	ClassId      string `json:"classId,omitempty"`
	ClassName    string `json:"className,omitempty"`
}

// ResIpPortMeta 内部数据对象-Ip端口对应的资源信息
type ResIpPortMeta struct {
	IpPort     types.IpPort `json:"ipPort"`
	CiId       string       `json:"ciId,omitempty"`
	CiName     string       `json:"ciName,omitempty"`
	CiTypeId   string       `json:"ciTypeId,omitempty"`
	CiTypeName string       `json:"ciTypeName,omitempty"`
}

// toResIpPortMeta 接口对象转内部对象
// @title       toResIpPortMeta
// @description 接口对象转内部对象
// @auth        小卒    2022/08/03 10:57
// @param       instMeta *InstanceMeta  "CMDB对象"
// @return      meta     *ResIpPortMeta "内部对象"
func toResIpPortMeta(instMeta *InstanceMeta) *ResIpPortMeta {
	meta := &ResIpPortMeta{
		IpPort: types.IpPort{
			Ip:   instMeta.Ip,
			Port: instMeta.Port,
		},
		CiId:       instMeta.InstanceId,
		CiName:     instMeta.InstanceName,
		CiTypeId:   instMeta.ClassId,
		CiTypeName: instMeta.ClassName,
	}
	return meta
}
