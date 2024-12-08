package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WaitOSSignalHandler(f func(), signals ...os.Signal) {
	if len(signals) == 0 {
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	<-ctx.Done()
	f()
}

func RegisterOSSignalHandler(f func(), signals ...os.Signal) {
	if len(signals) == 0 {
		return
	}

	go WaitOSSignalHandler(f, signals...)
}

func RegisterSignalInterruptHandler(f func()) {
	go WaitOSSignalHandler(f, os.Interrupt, syscall.SIGTERM)
}
