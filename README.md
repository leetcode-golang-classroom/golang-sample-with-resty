# golang-sample-with-resty

This repository is demo for how to use resty to send http request with easily way

## logic

```golang
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-resty/internal/service/task"
	"resty.dev/v3"
)

func main() {
  // structure logger
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		},
	))
	client := resty.New().SetBaseURL(config.AppConfig.ServerURI)
	// 1. Create a new task (POST)
	newTask := task.Task{Title: "Learn Resty", Done: false}
	var createdTask task.Task
	_, err := client.R().SetTimeout(2*time.Second).
		SetHeader("Content-Type", "application/json").
		SetBody(newTask).
		SetResult(&createdTask).
		Post("/tasks")
	if err != nil {
		logger.Error("failed to request post", slog.Any("err", err))
		os.Exit(1)
	}
	logger.Info("Created Task", slog.Any("task", createdTask))
	// 2. Get all tasks (GET)
	var tasks []task.Task
	_, err = client.R().
		SetResult(&tasks).
		Get("/tasks")
	if err != nil {
		logger.Error("Failed to get tasks", slog.Any("err", err))
		os.Exit(2)
	}
	logger.Info("All tasks", slog.Any("tasks", tasks))
	// 3. Update a task (PUT)
	updatedTask := task.Task{Title: "Master Resty", Done: true}
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(updatedTask).
		SetResult(&updatedTask).
		Put(fmt.Sprintf("/tasks/%d", createdTask.ID))
	if err != nil {
		logger.Error("Failed to update task", slog.Any("err", err))
		os.Exit(2)
	}
	logger.Info("updated task", slog.Any("task", updatedTask))
	// 4. Delete a task (DELETE)
	_, err = client.R().
		Delete(fmt.Sprintf("/tasks/%d", createdTask.ID))

	if err != nil {
		logger.Error("Failed to delete tasks", slog.Any("err", err))
		os.Exit(4)
	}
	logger.Info("Deleted Task", slog.Int("id", createdTask.ID))
}

```