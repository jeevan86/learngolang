package tcp

import (
	"fmt"
	logging "gopackettest/logger"
)

var logger = logging.LoggerFactory.NewLogger([]string{"stdout"}, []string{"stderr"})

type tcpChannel struct {
	srcIp, dstIp     string
	srcPort, dstPort int
}

func (o tcpChannel) ToString() string {
	return fmt.Sprintf("%s:%d->%s:%d", o.srcIp, o.srcPort, o.dstIp, o.dstPort)
}

// flag 是否启用助记符
func flag(enabled bool, name string) string {
	if enabled {
		return "[" + name + "]"
	}
	return ""
}
