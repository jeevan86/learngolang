package cancel

import (
	"context"
	"github.com/jeevan86/learngolang/pkg/util/panics"
)

type Cancellable struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	runnable   func()
}

func StartCancelable(runnable func()) *Cancellable {
	dr := NewCancelable(runnable)
	dr.Start()
	return dr
}

func NewCancelable(runnable func()) *Cancellable {
	// 新建一个上下文
	ctx := context.Background()
	// 在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancelFunc := context.WithCancel(ctx)
	return &Cancellable{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		runnable:   runnable,
	}
}

func (r *Cancellable) Start() {
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				break
			default:
				panics.SafeRun(r.runnable)
			}
		}
	}()
}

func (r *Cancellable) Cancel() {
	panics.SafeRun(r.cancelFunc)
}
