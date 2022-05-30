package util

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type FluxSink struct {
	input chan rxgo.Item
	flux  rxgo.Observable
}

func NewFluxSink(bufferSz int) *FluxSink {
	input := make(chan rxgo.Item, bufferSz)
	flux := rxgo.FromChannel(input) // rxgo.Just("Hello, World!")()
	fluxSink := FluxSink{input: input, flux: flux}
	return &fluxSink
}

func (sink *FluxSink) Subscribe(f func(e interface{})) {
	output := sink.flux.Observe()
	go func() {
		for item := range output {
			// do nothing
			f(item.V)
		}
	}()
}

func (sink *FluxSink) Next(e interface{}) {
	sink.input <- rxgo.Item{V: e}
}

func (sink *FluxSink) Map(f func(e interface{}) interface{}) *FluxSink {
	mp := func(ctx context.Context, o interface{}) (interface{}, error) {
		return f(o), nil
	}
	sink.flux.Map(mp)
	return sink
}

// FlatMap TODO: FlatMap实现起来比较麻烦
func (sink *FluxSink) FlatMap(f func(e interface{}) *FluxSink) *FluxSink {
	var flatMp func(rxgo.Item) rxgo.Observable
	var newFluxSink *FluxSink
	flatMp = func(item rxgo.Item) rxgo.Observable {
		//newFluxSink = f(item.V)
		// 根据item.V的值，创建一个数据新的流
		rxgo.Just(item.V)

		newFluxSink = NewFluxSink(2048)
		newFluxSink.Subscribe(func(e interface{}) {})
		newFluxSink.Next(item.V)

		output := newFluxSink.flux.Observe()
		return rxgo.FromChannel(output)
	}
	sink.flux.FlatMap(flatMp)
	return newFluxSink
}

func (sink *FluxSink) DoOnNext(f rxgo.NextFunc) *FluxSink {
	sink.flux.DoOnNext(f)
	return sink
}
