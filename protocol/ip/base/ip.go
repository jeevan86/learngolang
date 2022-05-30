package base

import (
	logging "gopackettest/logger"
)

var logger = logging.LoggerFactory.NewLogger([]string{"stdout"}, []string{"stderr"})

// Tags IP公共的选项
type Tags struct {
	srcIp, dstIp     string
	srcPort, dstPort int
}

// 包数量（增量）
type count uint64

// 字节数（增量）
type bytes uint64
