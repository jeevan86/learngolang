package cmdb

type serverConfig struct {
	url string `yaml:"url"`
}

type clientConfig struct {
}

type config struct {
	server serverConfig `yaml:"server"`
	client clientConfig `yaml:"client"`
}
