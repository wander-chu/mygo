package main

import (
	"context"
	"flag"
	"log/slog"
	"mygo/internal/option"
	"mygo/internal/presentation/httpapi"
	"mygo/pkg/logger"
	"os"
)

var (
	opt = &option.Options{}
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = httpapi.NewServer(opt)
}

func init() {
	flag.StringVar(&opt.ConfigFile, "configFile", "", "config file")
	flag.StringVar(&opt.LogLevel, "logLevel", "", "log level")
	flag.StringVar(&opt.DBDir, "dbDir", "", "database dir")
	flag.Parse()

	initLogger(opt)

	if opt.ConfigFile == "" {
		logAndExist("config file required")
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
