package actuator

// Metric 指标的结构
type Metric struct {
	Metric   string `json:"metric"`
	Instance string `json:"instance_name"`
	//"protocol": "tcp"
	//"protocol_ver": ""
	Tags      map[string]interface{} `json:"metric_tags"`
	Timestamp uint64                 `json:"timestamp"`
	Value     float64                `json:"value"`
}
