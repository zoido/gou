package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/zoido/gou/internal/exec/app/config"
	"github.com/zoido/gou/internal/exec/app/factory"
	"github.com/zoido/gou/internal/exec/exit"
	"github.com/zoido/gou/internal/exec/logging"
	"github.com/zoido/gou/internal/scanner"
	"github.com/zoido/gou/internal/ui"
)

var (
	shutdownTimeout   = 10 * time.Second
	hardShutdownDelay = 3 * time.Second
)

// Run sets up and start the GoU application.
func Run(ctx context.Context, args []string) error {
	ctx, stopHandlingSignals := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stopHandlingSignals()

	cfg, err := config.ParseConfig(args)
	if usageErr, shouldPrintUsage := errors.AsType[config.PrintUsageError](err); shouldPrintUsage {
		fmt.Println(usageErr.Usage)
		return nil
	}
	if err != nil {
		return exit.Newf(2, "failed to parse config: %v", err)
	}

	logging.Setup(cfg.Logging)
	slog.Debug("GoU config", "config", cfg)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GoU panic", "panic", r)
		}
	}()

	f := factory.New(cfg)

	open := f.InputOpener()
	reader, err := open()
	if err != nil {
		return exit.Newf(100, "opening input: %v", err)
	}
	closeInput := sync.OnceFunc(func() {
		if err := reader.Close(); err != nil {
			slog.Warn("Failed closing input", "error", err)
		}
	})

	grp, grpCtx := errgroup.WithContext(ctx)
	success := make(chan struct{})
	grp.Go(func() error {
		findings, err := scanner.FindURLs(grpCtx, reader)
		if err != nil {
			return err
		}
		closeInput() // We don't need it anymore so we do not need to hold it.

		pick, err := ui.Pick(grpCtx, findings)
		if err != nil {
			return err
		}
		fmt.Println(pick.URL())

		close(success)
		return nil
	})

	select {
	case <-grpCtx.Done():
		slog.Error("Error in execution, shutting down.")
	case <-ctx.Done():
		slog.Debug(fmt.Sprintf("%v, shutting down", context.Cause(ctx)))
	case <-success:
		// All good just shutdown.
	}

	stopHandlingSignals()

	go func() {
		closeInput()
	}()

	slog.Debug("Waiting for app to shutdown.")
	errCh := make(chan error)
	go func() {
		errCh <- grp.Wait()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return exit.Newf(128, "error: %v", err)
		}
	case <-time.After(shutdownTimeout):
		slog.Debug("Shutdown timed out, stopping forcefully.")
		slog.Debug("Waiting before exit", "delay", hardShutdownDelay)
		time.Sleep(hardShutdownDelay)
	}

	return nil
}
