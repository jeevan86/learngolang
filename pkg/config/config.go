package config

import (
	"fmt"
	lf4go "github.com/jeevan86/lf4go/factory"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"unsafe"
)

type config struct {
	Pcap      packetCapture       `yaml:"pcap"`
	Collector collector           `yaml:"collector"`
	Server    serverConfig        `yaml:"server"`
	Logging   lf4go.LoggingConfig `yaml:"logging"`
	NodeIp    string              `yaml:"node-ip"`
}

var configYml = "./config/config.yml"
var defaultYml = `
server:
  address: 127.0.0.1:8630
pcap:
  devices:
    - prefix: any
      duration: 2s
      snaplen: 120 # bytes
      promisc: true
  #    - prefix: ens192
  #      duration: 2s
  #      snaplen: 120 # bytes
  #      promisc: true
  routine:
    parallelism: 4
  reactor:
    buffer: 2048
collector:
  server-type: http # grpc | http | log
  server-addr: "http://127.0.0.1:8630/collect" # localhost:50051 | "http://127.0.0.1:8630/collect"
  parallelism: 1
  par-buff-size: 64
logging:
  factory: zap # zap | logrus
  formatter: normal # normal | json
  appenders:
    - type: stdout # file | stdout | ... (coming soon)
  root-name: learngolang
  root-level: INFO
  package-levels:
    "protocol/ip/tcp": WARN
`
var Config = loadConfigYml(configYml)

func loadConfigYml(path string) *config {
	if len(path) == 0 {
		path = configYml
	}
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("解析配置错误：%s，使用默认配置。%s", err.Error(), defaultYml))
		yml = *(*[]byte)(unsafe.Pointer(&defaultYml))
	}
	c := new(config)
	err = yaml.Unmarshal(yml, c)
	if err != nil {
		fmt.Println(fmt.Sprintf("解析配置错误：%s", err.Error()))
		return nil
	}
	c.NodeIp = getNodeIp()
	return c
}

// getNodeIp 先从环境变量NODE_IP_ADDR获得配置的IP，如果没有，再根据hostname获取
func getNodeIp() string {
	nodeIp, ok := os.LookupEnv("NODE_IP_ADDR")
	if !ok {
		hn, e := os.Hostname()
		if e == nil {
			ips, _ := net.LookupIP(hn)
			for _, ip := range ips {
				if !ip.IsLoopback() {
					nodeIp = ip.String()
					break
				}
			}
		}
	}
	return nodeIp
}
