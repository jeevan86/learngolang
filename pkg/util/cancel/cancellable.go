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

// StartCancelable
// @title       创建并启动一个可以取消的执行器
// @description 创建并启动一个可以取消的执行器
// @auth        小卒    2022/08/03 10:57
// @param       runnable func()       "处理函数"
// @return      r        *Cancellable "可以启动和取消"
func StartCancelable(runnable func()) *Cancellable {
	dr := NewCancelable(runnable)
	dr.Start()
	return dr
}

// NewCancelable
// @title       创建一个可以取消的执行器
// @description 创建一个可以取消的执行器（使用context.WithCancel）
// @auth        小卒    2022/08/03 10:57
// @param       runnable func()       "处理函数"
// @return      r        *Cancellable "可以启动和取消"
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

// Start
// @title       启动一个可以取消的执行器
// @description 启动一个可以取消的执行器（使用goroutine）
// @auth        小卒    2022/08/03 10:57
func (r *Cancellable) Start() {
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				break
			default:
				_, _ = panics.SafeRun(r.runnable)
			}
		}
	}()
}

// Cancel
// @title       取消执行器的执行
// @description 取消执行器的执行（调用context.CancelFunc）
// @auth        小卒    2022/08/03 10:57
func (r *Cancellable) Cancel() {
	_, _ = panics.SafeRun(r.cancelFunc)
}
