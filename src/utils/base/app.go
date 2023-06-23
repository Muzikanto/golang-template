package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type App struct {
	Modules   []IModule
	Engine    *gin.Engine
	Container dig.Container
}

func (r *App) Init() *App {
	gin.SetMode(gin.ReleaseMode)
	r.Engine = gin.New()

	return r
}
func (r *App) InitModules(globalContainer *dig.Container) *App {
	var container = globalContainer.Scope("APP_MODULE")

	//
	var appModule = Module{modules: r.Modules}

	appModule.Init(container, r.Engine)

	return r
}
func (r *App) AddModule(module interface{ IModule }) *App {
	r.Modules = append(r.Modules, module)

	return r
}
func (r *App) Listen(address string) error {
	fmt.Printf("[App] - Server listening at: %s\n\n", address)
	return r.Engine.Run(address)
}
