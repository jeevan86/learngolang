package parallel

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/util/cancel"
)

// Routine 并发执行器
type Routine struct {
	ch      chan interface{}
	runners []*cancel.Cancellable
}

// Routines 并发调度入口对象
type Routines struct {
	parallelism uint64
	routines    []*Routine
}

// NewParRoutines
// @title       创建一个并发goroutine池
// @description 根据并发度创建协程
// @auth        小卒    2022/08/03 10:57
// @param       parallelism int               "并发度（几个goroutine）"
// @param       chBufSz     int               "每个goroutine维护的chan的缓存大小"
// @param       process     func(interface{}) "处理函数"
// @return      r           *Routines         "并发调度入口对象"
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

// adjustParallelism
// @title       调整并发参数
// @description 调整并发参数，并发度的合理范围[1~16]
// @auth        小卒        2022/08/03 10:57
// @param       parallelism int  "并发度（几个goroutine）"
// @return      r           int  "并发度"
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

// Dispatch
// @title       根据id将数据送到并发goroutine处理
// @description 根据id将数据送到并发goroutine处理
// @auth        小卒 2022/08/03 10:57
// @param       idx uint64      "数据的编号"
// @param       o   interface{} "数据对象"
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
