package config

const (
	ParTypeRoutine = "routine"
	ParTypeReactor = "reactor"
)

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

type routineConfig struct {
	Parallelism int `yaml:"parallelism"`
}

type reactorConfig struct {
	BufferSz int `yaml:"buffer"`
}
