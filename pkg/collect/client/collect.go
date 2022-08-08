package client

import (
	"github.com/jeevan86/learngolang/pkg/collect/api"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/util/panics"
	"github.com/jeevan86/learngolang/pkg/util/parallel"
	"strings"
)

var logger = log.NewLogger()

var collectConfig = config.GetConfig().Agent.Collect

type Collect struct {
	Api      api.CollectorApi
	input    chan *base.OutputStruct
	routines *parallel.Routines
}

const (
	TypeLog  = "log"
	TypeGrpc = "grpc"
	TypeHttp = "http"
)

// NewCollector
// @title       创建数据收集器
// @description 创建数据收集器，有三种方式：通过GRPC发送、通过HTTP发送，或者写本地日志
// @auth        小卒   2022/08/03 10:57
func NewCollector() *Collect {
	var collector api.CollectorApi
	serverType := *collectConfig.ServerType
	serverAddr := *collectConfig.ServerAddr
	switch strings.ToLower(serverType) {
	case TypeGrpc:
		collector = newGrpcCollector(serverAddr)
		break
	case TypeHttp:
		collector = newHttpCollector(serverAddr)
		break
	case TypeLog:
		collector = newLogCollector("log-collector")
		break
	default:
		collector = newLogCollector("log-collector")
		break
	}
	return &Collect{
		Api: collector,
	}
}

// Start
// @title       启动数据收集
// @description 启动数据收集
// @auth        小卒     2022/08/03 10:57
// @param       input chan *base.OutputStruct "数据输入源"
func (c *Collect) Start(input chan *base.OutputStruct) {
	c.input = input
	parallelism := *collectConfig.Parallelism
	parBuffSize := *collectConfig.ParBuffSize
	parFunction := c.makeParallel(parallelism, parBuffSize)
	go func() {
		idx, after := cyclingIdxFunc()
		for {
			m, ok := <-c.input
			if !ok {
				break
			}
			parFunction(*idx, m) // single routine use c.routine(m)
			after()
		}
	}()
}

// cyclingIdxFunc
// @title       可循环的ID、和递增函数
// @description 可循环的ID、和递增函数，不支持并发。
// @auth        小卒    2022/08/03 10:57
// @return      idx   *uint64  "IP端口的信息"
// @return      after func()   "递增函数"
func cyclingIdxFunc() (*uint64, func()) {
	// 一次循环使用的ID的最大值
	const maxIdxPerCycle = ^uint64(0) - 99999999
	// 表示第几次循环
	cycle := uint64(0)
	// 当前ID值
	idx := uint64(0)
	// id函数递增
	after := func() {
		if idx >= maxIdxPerCycle {
			idx = 0
			cycle++
		}
		idx++
	}
	return &idx, after
}

// makeParallel
// @title       并发处理（多协程处理）
// @description 根据并发度创建协程
// @auth        小卒    2022/08/03 10:57
// @param       parallelism int  "并发度（几个goroutine）"
// @param       chBufSz     int  "每个goroutine维护的chan的缓存地址"
// @return      f           func(uint64, *base.OutputStruct) "接口函数"
func (c *Collect) makeParallel(parallelism, chBufSz int) func(uint64, *base.OutputStruct) {
	// with Routines
	c.routines = parallel.NewParRoutines(
		parallelism,
		chBufSz,
		false,
		func(o interface{}) {
			c.routine(o.(*base.OutputStruct))
		},
	)
	return func(i uint64, o *base.OutputStruct) {
		sta, err := panics.SafeRun(
			func() {
				c.routines.Dispatch(i, o)
			},
		)
		if err != nil {
			logger.Error("Failed to dispatch data, err => %s, stack => %s", err.Error(), sta)
		}
	}
}

// routine
// @title       并发协程处理逻辑
// @description 并发协程处理逻辑
// @auth        小卒    2022/08/03 10:57
// @param       m  *base.OutputStruct  "并发协程处理的输入消息"
func (c *Collect) routine(m *base.OutputStruct) {
	sta, err := panics.SafeRun(
		func() {
			c.Api.Collect(m)
		},
	)
	if err != nil {
		logger.Error("Failed to collect data, err => %s, stack => %s", err.Error(), sta)
	}
}

func (c *Collect) Stop() {
	c.routines.Close()
	close(c.input)
}
