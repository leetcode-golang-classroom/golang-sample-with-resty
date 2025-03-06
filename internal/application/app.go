package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/logger"
)

type App struct {
	appConfig *config.Config
	router    *http.ServeMux
}

// New - 建立 App 物件
func New(ctx context.Context, appConfig *config.Config) *App {
	router := http.NewServeMux()
	// create app instance
	app := &App{
		appConfig: appConfig,
		router:    router,
	}
	// setup route
	app.setupTaskRoute(ctx)
	return app
}

// Start - 啟動 server
func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.appConfig.Port),
		Handler: app.router,
	}
	log := logger.FromContext(ctx)
	log.Info(fmt.Sprintf("starting server on %s", app.appConfig.Port))
	errCh := make(chan error, 1)
	var err error
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
		CloseChannel(errCh)
	}()
	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		log.Warn("stopping server, wait for 10 seconds to stop")
		defer cancel()
		return server.Shutdown(timeout)
	}
}

func CloseChannel(ch chan error) {
	if _, ok := <-ch; ok {
		close(ch)
	}
}
