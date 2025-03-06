package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/application"
	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/config"
	mlog "github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/logger"
)

func main() {
	// 建立 logger
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		},
	))
	rootContext := context.WithValue(context.Background(), mlog.CtxKey{}, logger)
	// 建立 application instance
	app := application.New(rootContext, config.AppConfig)
	// 設定中斷訊號監聽
	ctx, cancel := signal.NotifyContext(rootContext, os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)

	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		logger.Error("failed to start app", slog.Any("err", err))
	}
}
