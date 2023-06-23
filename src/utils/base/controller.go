package base

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type IController interface {
	Init(container *dig.Scope, engine *gin.Engine)
}

type Controller struct {
	IController
	Engine    *gin.Engine
	Container *dig.Scope
}

func (r Controller) Init(container *dig.Scope, engine *gin.Engine) {
	r.Engine = engine
	r.Container = container
}
