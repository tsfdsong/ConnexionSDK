package skywalking

import (
	"context"

	"github.com/SkyAPM/go2sky"
)

type StepHandler interface {
	Handle(ctx context.Context, seq int) error
	GetContext() context.Context
	GetSpan() go2sky.Span
	End()
}

type Sequence struct {
	ctx      context.Context
	handlers []StepHandler
	seq      int
}

func (s *Sequence) Init(c context.Context) {
	s.ctx = c
	s.seq = 1
}

func (s *Sequence) SetContext(c context.Context) {
	s.ctx = c
}

func (s *Sequence) AddStep(h StepHandler) {
	s.handlers = append(s.handlers, h)
}

func (s *Sequence) Execute() error {
	for _, v := range s.handlers {
		err := v.Handle(s.ctx, s.seq)
		if err != nil {
			return err
		}
		s.SetContext(v.GetContext())
		s.seq = s.seq + 1

	}
	return nil
}

func (s *Sequence) GetStepHandler(i int) StepHandler {
	return s.handlers[i]
}

func (s *Sequence) End() {
	for _, v := range s.handlers {
		v.End()
	}
}

func (s *Sequence) Iterator() func() (StepHandler, bool) {
	index := 0
	return func() (val StepHandler, ok bool) {
		if index >= len(s.handlers) {
			return
		}

		val, ok = s.handlers[index], true
		index++
		return
	}
}
