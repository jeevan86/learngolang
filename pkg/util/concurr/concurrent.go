package concurr

import (
	"github.com/jeevan86/learngolang/pkg/util/cancel"
	"github.com/jeevan86/learngolang/pkg/util/panics"
)

type Sync struct {
	disposable *cancel.Cancellable
	ch         chan *SyncCommand
}

type SyncCommand struct {
	Var []interface{}
	Fun func(...interface{})
}

func (s *Sync) SyncProcess(cmd *SyncCommand) {
	s.ch <- cmd
}

func (s *Sync) Start() {
	f := func(cmd *SyncCommand) {
		_, _ = panics.SafeRun(func() { cmd.Fun(cmd.Var...) })
	}
	s.disposable = cancel.NewCancelable(func() {
		for {
			c, ok := <-s.ch
			if !ok {
				break
			}
			f(c)
		}
	})
	s.disposable.Start()
}

func (s *Sync) Stop() {
	if s.disposable != nil {
		s.disposable.Cancel()
	}
}
