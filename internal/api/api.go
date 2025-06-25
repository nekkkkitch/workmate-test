package api

import "github.com/gofiber/fiber/v2"

type API struct {
	App     *fiber.App
	config  *Config
	service Service
}

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func New(cfg Config, service Service) (*API, error) {
	app := fiber.New()
	api := API{App: app, config: &cfg, service: service}
	app.Post("/task", api.CreateTask)
	app.Get("/tasks", api.GetAllTasks)
	app.Get("/task/+", api.GetTask)
	app.Delete("/task/+", api.DeleteTask)
	app.Listen(cfg.Host + ":" + cfg.Port)
	return &api, nil
}
