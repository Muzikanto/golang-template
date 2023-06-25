package core

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/engine"
	"go-backend-template/src/utils/logger"
	"go.uber.org/dig"
	"log"
)

type Controller struct {
	Container *dig.Scope
	Engine    *engine.Engine
	Router    *gin.RouterGroup
	Logger    *logger.Logger
}

func CreateController(container *dig.Scope, engine *engine.Engine, path string) *Controller {
	var controller = &Controller{Container: container, Engine: engine, Router: engine.Group(path)}

	if err := container.Invoke(func(logger *logger.Logger) {
		controller.Logger = logger
	}); err != nil {
		log.Fatal(err)
	}

	return controller
}

func (receiver *Controller) BindBody(payload interface{}, c *gin.Context) error {
	return BindBody(payload, c)
}

func (receiver *Controller) ErrorResponse(err error, data interface{}, withDetails bool) *Response {
	return ErrorResponse(err, data, withDetails)
}

func (receiver *Controller) OkResponse(data interface{}) *Response {
	return OkResponse(data)
}
