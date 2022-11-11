package flow

import (
	"context"
	"fmt"
	"time"
)

type Step struct {
	name string
	fn   func(context.Context) error
}

type Flow struct {
	ctx   context.Context
	steps []Step
}

func (f *Flow) Step(name string, fn func(context.Context) error) *Flow {
	f.steps = append(f.steps, Step{name: name, fn: fn})
	return f
}

func (f *Flow) Run() {
	start := time.Now()
	for _, v := range f.steps {
		if err := v.fn(f.ctx); err != nil {
			fmt.Printf("[x]run %s failed, err = %v\n", v.name, err)
			break
		}
		fmt.Printf("[o]run %s done\n", v.name)
	}
	fmt.Printf("[o]flow took %.2fs\n", time.Since(start).Seconds())
}

func NewFlow(ctx context.Context) *Flow {
	return &Flow{ctx: ctx}
}
