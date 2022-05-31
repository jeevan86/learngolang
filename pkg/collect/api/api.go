package api

import "github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"

type CollectorApi interface {
	Collect(*base.OutputStruct)
	GetLocalIpList(string) []string
}
