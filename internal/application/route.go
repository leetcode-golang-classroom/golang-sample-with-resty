package application

import (
	"context"

	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/logger"
	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/service/task"
)

func (app *App) setupTaskRoute(ctx context.Context) {
	taskRoute := task.NewTaskHandler(logger.FromContext(ctx))
	taskRoute.SetupRoute(app.router)
}
