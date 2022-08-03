package backend

type IpPortType string

const (
	IpPortTypeK8sNode   IpPortType = "k8sNode"
	IpPortTypeK8sPod    IpPortType = "k8sPod"
	IpPortTypeK8sSrv    IpPortType = "k8sService"
	IpPortTypeComponent IpPortType = "component"
	IpPortTypeUnknown   IpPortType = "unknown"
)

type IpFamily int8

const (
	IpFamilyIp4 IpFamily = 4
	IpFamilyIp6 IpFamily = 6
)

type IpProtocol string

const (
	IpProtocolTcp4 IpProtocol = "tcp4"
	IpProtocolUdp4 IpProtocol = "udp4"
	IpProtocolTcp6 IpProtocol = "tcp6"
	IpProtocolUdp6 IpProtocol = "udp6"
)

type Tags map[string]string

type IpPortMeta struct {
	Type IpPortType `json:"type,omitempty"`
	Tags Tags       `json:"tags,omitempty"`
}

type PacketType string

const (
	PacketTypeReq     PacketType = "req"
	PacketTypeRes     PacketType = "res"
	PacketTypeUnknown PacketType = "unknown"
)

type ChannelPacketMeta struct {
	Protocol IpProtocol  `json:"protocol,omitempty"`
	Type     PacketType  `json:"type,omitempty"`
	SrcIp    string      `json:"srcIp,omitempty"`
	SrcPort  int32       `json:"srcPort,omitempty"`
	SrcMeta  *IpPortMeta `json:"srcMeta,omitempty"`
	DstIp    string      `json:"dstIp,omitempty"`
	DstPort  int32       `json:"dstPort,omitempty"`
	DstMeta  *IpPortMeta `json:"dstMeta,omitempty"`
}
