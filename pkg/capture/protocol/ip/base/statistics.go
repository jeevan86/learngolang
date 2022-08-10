package base

// CommonAggregatedValues 通用的ip包聚合数据结构
type CommonAggregatedValues struct {
	SendCount uint64 `json:"send_count"`
	SendBytes uint64 `json:"send_bytes"`
	RecvCount uint64 `json:"recv_count"`
	RecvBytes uint64 `json:"recv_bytes"`
	Count     uint64 `json:"count"`
	Bytes     uint64 `json:"bytes"`
}
