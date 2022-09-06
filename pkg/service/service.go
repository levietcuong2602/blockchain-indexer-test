package service

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunWithGracefulShutdown(ctx context.Context, runFunc func(ctx context.Context, wg *sync.WaitGroup)) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	runFunc(ctx, wg)

	<-stop

	cancel()
	wg.Wait()
}
