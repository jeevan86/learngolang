package config

type agentConfig struct {
	Pcap    *packetCapture `yaml:"pcap"`
	Collect *collect       `yaml:"collect"`
}

type collect struct {
	ServerType  *string `yaml:"server-type"`
	ServerAddr  *string `yaml:"server-addr"`
	Parallelism *int    `yaml:"parallelism"`
	ParBuffSize *int    `yaml:"par-buff-size"`
}

type packetCapture struct {
	Devices []deviceConfig `yaml:"devices"`
	ParType string         `yaml:"par-type"`
	Routine routineConfig  `yaml:"routine"`
	Reactor reactorConfig  `yaml:"reactor"`
}

type deviceConfig struct {
	Prefix   string `yaml:"prefix"`
	Duration string `yaml:"duration"`
	Snaplen  int32  `yaml:"snaplen"`
	Promisc  bool   `yaml:"promisc"`
}

const (
	ParTypeRoutine = "routine"
	ParTypeReactor = "reactor"
)

type routineConfig struct {
	Parallelism int `yaml:"parallelism"`
}

type reactorConfig struct {
	BufferSz int `yaml:"buffer"`
}
