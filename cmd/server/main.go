package main

import (
	"context"
	"flag"
	"log/slog"
	"mygo/internal/option"
	"mygo/internal/presentation/httpapi"
	"mygo/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	opt = &option.Options{}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := httpapi.NewServer(opt)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		logger.Debug(ctx, "receive signal", "signal", s)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		if err := server.Close(wg); err != nil {
			logger.Error(ctx, "shutdown server", "error", err)
		} else {
			logger.Info(ctx, "shutdown server")
		}

		wg.Wait()
		os.Exit(0)
	}
}

func init() {
	flag.StringVar(&opt.ConfigFile, "configFile", "", "config file")
	flag.StringVar(&opt.LogLevel, "logLevel", "", "log level")
	flag.StringVar(&opt.DBDir, "dbDir", "", "database dir")
	flag.Parse()

	initLogger(opt)

	if opt.ConfigFile == "" {
		logAndExist("config file required")
	} else if err := opt.LoadFile(opt.ConfigFile); err != nil {
		logAndExist("load config file", "error", err)
	} else if err := opt.Prepare(); err != nil {
		logAndExist("prepare resources", "error", err)
	}
}

func initLogger(opt *option.Options) {
	levels := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	level := slog.LevelInfo
	if v, ok := levels[opt.LogLevel]; ok {
		level = v
	}

	slog.SetDefault(slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}),
	))
}

func logAndExist(message string, args ...any) {
	logger.Error(context.TODO(), message, args...)
	os.Exit(1)
}
