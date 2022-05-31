package base

type CommonAggregatedValues struct {
	SendCount uint64 `yaml:"send-count" json:"send_count"`
	SendBytes uint64 `yaml:"send-bytes" json:"send_bytes"`
	RecvCount uint64 `yaml:"recv-count" json:"recv_count"`
	RecvBytes uint64 `yaml:"recv-bytes" json:"recv_bytes"`
	Count     uint64 `yaml:"count" json:"count"`
	Bytes     uint64 `yaml:"bytes" json:"bytes"`
}
