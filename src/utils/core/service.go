package core

import (
	"go-backend-template/src/utils/logger"
	"go.uber.org/dig"
	"log"
)

type Service struct {
	Container *dig.Scope
	Logger    *logger.Logger
}

func CreateService(container *dig.Scope) *Service {
	var service = &Service{Container: container}

	if err := container.Invoke(func(logger *logger.Logger) {
		service.Logger = logger
	}); err != nil {
		log.Fatal(err)
	}

	return service
}
