package config

type agentConfig struct {
	Capture *captureConfig `yaml:"capture"`
	Collect *collectConfig `yaml:"collect"`
}

type collectConfig struct {
	ServerType  *string `yaml:"server-type"`
	ServerAddr  *string `yaml:"server-addr"`
	Parallelism *int    `yaml:"parallelism"`
	ParBuffSize *int    `yaml:"par-buff-size"`
}

type captureConfig struct {
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
	Parallelism  int  `yaml:"parallelism"`
	ChBufferSize int  `yaml:"ch-buffer-size"`
	ShareChan    bool `yaml:"share-chan"`
}

type reactorConfig struct {
	BufferSz int `yaml:"buffer"`
}
