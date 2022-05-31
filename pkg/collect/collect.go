package collect

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

var collectorConfig = config.Config.Collector

type Collector struct {
	Api      api.CollectorApi
	input    chan *base.OutputStruct
	routines *parallel.Routines
}

const (
	TypeLog  = "log"
	TypeGrpc = "grpc"
	TypeHttp = "http"
)

func NewCollector() *Collector {
	var collector api.CollectorApi
	switch strings.ToLower(collectorConfig.ServerType) {
	case TypeGrpc:
		collector = newGrpcCollector(collectorConfig.ServerAddr)
		break
	case TypeHttp:
		//collector = newHttpCollector(collectorConfig.ServerAddr)
		collector = newHttpCollector(collectorConfig.ServerAddr)
		break
	case TypeLog:
		collector = newLogCollector("log-collector")
		break
	default:
		collector = newLogCollector("log-collector")
		break
	}
	return &Collector{
		Api: collector,
	}
}

func (c *Collector) Start(input chan *base.OutputStruct) {
	c.input = input
	parallelism := collectorConfig.Parallelism
	parBuffSize := collectorConfig.ParBuffSize
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

func cyclingIdxFunc() (*uint64, func()) {
	const maxIdxPerCycle = ^uint64(0) - 99999999
	cycle := uint64(0)
	idx := uint64(0)
	after := func() {
		if idx >= maxIdxPerCycle {
			idx = 0
			cycle++
		}
		idx++
	}
	return &idx, after
}

func (c *Collector) makeParallel(parallelism, chBufSz int) func(uint64, *base.OutputStruct) {
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

func (c *Collector) routine(m *base.OutputStruct) {
	sta, err := panics.SafeRun(
		func() {
			c.Api.Collect(m)
		},
	)
	if err != nil {
		logger.Error("Failed to collect data, err => %s, stack => %s", err.Error(), sta)
	}
}

func (c *Collector) Stop() {
	c.routines.Close()
	close(c.input)
}
