package types

type IpPortType string

const (
	IpPortTypeK8sNode   IpPortType = "k8sNode"
	IpPortTypeK8sPod    IpPortType = "k8sPod"
	IpPortTypeK8sSrv    IpPortType = "k8sService"
	IpPortTypeComponent IpPortType = "component"
	IpPortTypeHost      IpPortType = "host"
	IpPortTypeUnknown   IpPortType = "unknown"
)

type IpFamily int8

const (
	IpFamilyIp4 IpFamily = 4
	IpFamilyIp6 IpFamily = 6
)

type IpProtocol string

const (
	IpProtocolTcp4    IpProtocol = "tcp4"
	IpProtocolUdp4    IpProtocol = "udp4"
	IpProtocolTcp6    IpProtocol = "tcp6"
	IpProtocolUdp6    IpProtocol = "udp6"
	IpProtocolUnknown IpProtocol = "unknown"
)

type Tags map[string]string

type PacketType string

const (
	PacketTypeReq     PacketType = "req"
	PacketTypeRes     PacketType = "res"
	PacketTypeUnknown PacketType = "unknown"
)

type IpPort struct {
	Ip   string `json:"ip,omitempty" yaml:"ip"`
	Port int32  `json:"port,omitempty" yaml:"port"`
}

type IpPortMeta struct {
	IpPort
	Type IpPortType `json:"type,omitempty"`
	Tags Tags       `json:"tags,omitempty"`
}

type ChannelPacketMeta struct {
	Protocol  IpProtocol  `json:"protocol,omitempty"`
	Type      PacketType  `json:"type,omitempty"`
	SrcIpPort *IpPortMeta `json:"srcIpPort,omitempty"`
	DstIpPort *IpPortMeta `json:"dstIpPort,omitempty"`
}
