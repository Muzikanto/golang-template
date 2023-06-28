package core

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/engine"
	error2 "go-backend-template/src/utils/error"
	"go-backend-template/src/utils/logger"
	"go.uber.org/dig"
	"log"
	"net/http"
)

type Controller struct {
	Container *dig.Scope
	Engine    *engine.Engine
	Router    *gin.RouterGroup
	Logger    *logger.Logger
}

func CreateController(container *dig.Scope, engine *engine.Engine, path string) *Controller {
	var router = engine.Group(path)
	router.Use(jsonLoggerMiddleware())

	var controller = &Controller{Container: container, Engine: engine, Router: router}

	if err := container.Invoke(func(logger *logger.Logger) {
		controller.Logger = logger
	}); err != nil {
		log.Fatal(err)
	}

	return controller
}

func (r *Controller) BindBody(payload interface{}, c *gin.Context) error {
	return BindBody(payload, c)
}

func (r *Controller) BindQuery(payload interface{}, c *gin.Context) error {
	return BindQuery(payload, c)
}

func (r *Controller) Bind(payload interface{}, c *gin.Context) error {
	return Bind(payload, c)
}

func (r *Controller) ErrorResponse(err error, data interface{}, withDetails bool) *Response {
	return ErrorResponse(err, data, withDetails)
}

func (r *Controller) OkResponse(data interface{}) *Response {
	return OkResponse(data)
}

// helpers

func (r *Controller) Authorize(c *gin.Context) {
	var token, err = c.Cookie("token")

	if err != nil {
		var err = error2.New(error2.UnauthorizedError, "Not Authorized")

		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err, nil, false))

		return
	}

	var user interface{}

	c.Set("token", token)
	c.Set("user", user)

	c.Next()
}

// guards

func (r *Controller) AuthGuard(c *gin.Context) {
	var user, exists = c.Get("User")

	if !exists || user == nil {
		var err = error2.New(error2.UnauthorizedError, "Not Authorized")

		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err, nil, false))

		return
	}

	c.Next()
}

func (r *Controller) RoleGuard(roles []string) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var roleInterface, exists = c.Get("Role")

		if !exists {
			if role, ok := roleInterface.(string); ok {
				for _, a := range roles {
					if a == role {
						c.Next()
						return
					}
				}
			}
		}

		var err = error2.New(error2.UnauthorizedError, "Permission Denied")

		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err, nil, false))
	}
}
