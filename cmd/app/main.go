package main

import (
	"log/slog"
	"workmate/internal/api"
	"workmate/internal/service"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServiceConfig *service.Config `yaml:"service"`
	APIConfig     *api.Config     `yaml:"api"`
}

func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := readConfig("./cfg.yml")
	if err != nil {
		slog.Error("Can't read config", "error", err)
		return
	}
	slog.Info("Config read successfully", "config", cfg)
	svc, err := service.New(*cfg.ServiceConfig)
	if err != nil {
		slog.Error("can't create service", "error", err)
		return
	}
	go svc.Start()
	slog.Info("Service created successfully")
	_, err = api.New(*cfg.APIConfig, svc)
	if err != nil {
		slog.Error("can't create api", "error", err)
		return
	}
	slog.Info("api created successfully")
}
