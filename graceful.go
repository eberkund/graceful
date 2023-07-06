package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var ErrSignalExit = errors.New("exit signal received")

type TaskFunc func() error

// Graceful implements graceful cleanup for multiple subtasks.
type Graceful struct {
	ctx     context.Context
	cancel  context.CancelCauseFunc
	cleanup []TaskFunc
	wg      sync.WaitGroup
}

// New returns an instance of Graceful.
func New(ctx context.Context) (*Graceful, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	g := &Graceful{
		ctx:    ctx,
		cancel: cancel,
	}
	go g.signals()
	return g, ctx
}

func (g *Graceful) signals() {
	defer g.cancel(ErrSignalExit)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}

func (g *Graceful) shutdown() error {
	var (
		wg   sync.WaitGroup
		lock sync.Mutex
		errs []error
	)
	wg.Add(len(g.cleanup))
	for _, fn := range g.cleanup {
		go func(fn TaskFunc) {
			defer wg.Done()
			err := fn()
			lock.Lock()
			errs = append(errs, err)
			lock.Unlock()
		}(fn)
	}
	wg.Wait()
	return errors.Join(errs...)
}

// Go adds a function to the group.
func (g *Graceful) Go(fn TaskFunc) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		err := fn()
		g.cancel(err)
	}()
}

// Stop adds a function that will be executed when group stops.
func (g *Graceful) Stop(fn TaskFunc) {
	g.cleanup = append(g.cleanup, fn)
}

// WaitWithErrors blocks until all functions are done and returns
// the shutdown cause and any errors from the cleanup functions.
func (g *Graceful) WaitWithErrors() (error, error) {
	<-g.ctx.Done()
	err := g.shutdown()
	g.wg.Wait()
	cause := context.Cause(g.ctx)
	return cause, err
}

// Wait blocks until all functions are done and returns the shutdown cause.
func (g *Graceful) Wait() error {
	err, _ := g.WaitWithErrors()
	return err
}
