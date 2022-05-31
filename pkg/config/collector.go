package config

type collector struct {
	ServerType  string `yaml:"server-type"`
	ServerAddr  string `yaml:"server-addr"`
	Parallelism int    `yaml:"parallelism"`
	ParBuffSize int    `yaml:"par-buff-size"`
}
