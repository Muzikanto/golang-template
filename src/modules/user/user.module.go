package user

import (
	"go-backend-template/src/utils/engine"
	"go-backend-template/src/utils/logger"
	"go.uber.org/dig"
	"log"
)

type Module struct {
	Container *dig.Scope
	Logger    *logger.Logger
}

func CreateModule(parentContainer *dig.Scope, engine *engine.Engine) *Module {
	var module = &Module{}
	module.Container = parentContainer.Scope("UserModule")

	// inject
	if err := module.Container.Invoke(func(logger *logger.Logger) {
		module.Logger = logger.Clone("User").Provide(module.Container)
	}); err != nil {
		log.Fatal(err)
	}

	// provide
	if err := module.Container.Provide(func() *UserService {
		return CreateTestService(module.Container)
	}); err != nil {
		log.Fatal(err)
	}

	var _ = CreateController(module.Container, engine)

	//
	module.Logger.Log("UserModule initialized")

	return module
}
