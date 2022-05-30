package util

import "context"

type DisposableGoRoutine struct {
	ctx      context.Context
	cancel   context.CancelFunc
	runnable func()
}

func StartDisposable(runnable func()) *DisposableGoRoutine {
	dr := newDisposableGoRoutine(runnable)
	dr.start()
	return dr
}

func newDisposableGoRoutine(runnable func()) *DisposableGoRoutine {
	// 新建一个上下文
	ctx := context.Background()
	// 在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)
	return &DisposableGoRoutine{
		ctx:      ctx,
		cancel:   cancel,
		runnable: runnable,
	}
}

func (r *DisposableGoRoutine) Dispose() {
	r.cancel()
}

func (r *DisposableGoRoutine) start() {
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				// zap.Error()
				break
			default:
				r.runnable()
			}
		}
	}()
}
