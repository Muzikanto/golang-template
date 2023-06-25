package engine

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"log"
)

type Engine struct {
	*gin.Engine
}

func CreateEngine() *Engine {
	gin.SetMode(gin.ReleaseMode)

	return &Engine{Engine: gin.New()}
}

func (r *Engine) Provide(container *dig.Container) *Engine {
	if err := container.Provide(func() *Engine {
		return r
	}); err != nil {
		log.Fatal(err)
	}

	return r
}
