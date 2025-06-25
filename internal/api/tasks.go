package api

import (
	"log/slog"
	cerr "workmate/pkg/customErrors"
	"workmate/pkg/task"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Service interface {
	CreateTask() error
	GetAllTasks() ([]task.Task, error)
	GetTask(uuid.UUID) (task.Task, error)
	DeleteTask(uuid.UUID) error
}

func (a *API) CreateTask(c *fiber.Ctx) error {
	slog.Info("API: got CreateTask")
	err := a.service.CreateTask()
	if err != nil {
		return fiber.NewError(500, "Error: "+err.Error())
	}
	return c.SendStatus(200)
}

func (a *API) GetAllTasks(c *fiber.Ctx) error {
	slog.Info("API: got GetAllTasks")
	tasks, err := a.service.GetAllTasks()
	if err != nil {
		return fiber.NewError(500, "Error: "+err.Error())
	}
	return c.JSON(tasks)
}

func (a *API) GetTask(c *fiber.Ctx) error {
	slog.Info("API: got GetTask")
	taskID, err := uuid.Parse(c.Params("+"))
	if err != nil {
		return fiber.NewError(500, "Error: "+err.Error())
	}
	task, err := a.service.GetTask(taskID)
	if err != nil {
		if err == cerr.NoSuchTask {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return fiber.NewError(500, "Error: "+err.Error())
	}
	return c.JSON(task)
}

func (a *API) DeleteTask(c *fiber.Ctx) error {
	slog.Info("API: got DeleteTask")
	taskID, err := uuid.Parse(c.Params("+"))
	if err != nil {
		return fiber.NewError(500, "Error: "+err.Error())
	}
	err = a.service.DeleteTask(taskID)
	if err != nil {
		if err == cerr.NoSuchTask {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return fiber.NewError(500, "Error: "+err.Error())
	}
	return c.SendStatus(200)
}
