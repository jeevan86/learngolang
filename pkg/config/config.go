package config

import (
	flag2 "github.com/jeevan86/learngolang/pkg/flag"
	lf4go "github.com/jeevan86/lf4go/factory"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"unsafe"
)

type config struct {
	Pcap    packetCapture       `yaml:"pcap"`
	Collect collect             `yaml:"collect"`
	Server  serverConfig        `yaml:"server"`
	Logging lf4go.LoggingConfig `yaml:"logging"`
	NodeIp  string              `yaml:"node-ip"`
}

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
collect:
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
var instance *config
var mutex = sync.Mutex{}

func GetConfig() *config {
	if instance != nil {
		return instance
	} else {
		load()
	}
	return instance
}

func load() {
	mutex.Lock()
	defer mutex.Unlock()
	if instance != nil {
		return
	}
	path := *flag2.ConfFile
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("解析配置错误：%s，使用默认配置。%s", err.Error(), defaultYml)
		yml = *(*[]byte)(unsafe.Pointer(&defaultYml))
	}
	c := new(config)
	err = yaml.Unmarshal(yml, c)
	if err != nil {
		log.Fatalf("解析配置错误：%s", err.Error())
	}
	c.NodeIp = getNodeIp()
	instance = c
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
