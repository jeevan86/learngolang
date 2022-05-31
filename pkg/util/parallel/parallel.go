package parallel

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/util/cancel"
)

type Routine struct {
	ch      chan interface{}
	runners []*cancel.Cancellable
}

type Routines struct {
	parallelism uint64
	routines    []*Routine
}

func NewParRoutines(parallelism int, chBuffSz int, shareCh bool, process func(interface{})) *Routines {
	par := adjustParallelism(parallelism)
	internal := make([]*Routine, par)
	routines := Routines{
		parallelism: uint64(par),
		routines:    internal,
	}
	// 一般多个goroutine使用同一个channel，就可以了
	var ch chan interface{}
	if shareCh {
		ch = make(chan interface{}, chBuffSz)
	}
	// 初始化几个协程来处理数据
	for i := 0; i < par; i++ {
		routine := Routine{}
		if shareCh {
			routine.ch = ch
		} else {
			routine.ch = make(chan interface{}, chBuffSz)
		}
		routine.runners = []*cancel.Cancellable{
			cancel.StartCancelable(func() {
				o, ok := <-routine.ch
				if ok {
					process(o)
				}
			}),
		}
		routines.routines[i] = &routine
	}
	return &routines
}

func adjustParallelism(parallelism int) int {
	par := parallelism
	if parallelism < 1 || parallelism > 16 {
		fmt.Printf(
			"Impossible parallelism value: %d, range[1~16], use default 2",
			parallelism,
		)
		par = 2
	}
	return par
}

func (r *Routines) Dispatch(i uint64, o interface{}) {
	routines := r.routines
	parallelism := r.parallelism
	routines[i%parallelism].ch <- o
}

func (r *Routines) Close() {
	for _, routine := range r.routines {
		for _, disposable := range routine.runners {
			disposable.Cancel()
		}
	}
}
