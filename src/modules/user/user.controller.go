package user

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/core"
	"go-backend-template/src/utils/engine"
	"go.uber.org/dig"
	"log"
)

type UserController struct {
	*core.Controller
	testService *UserService
}

func CreateController(container *dig.Scope, engine *engine.Engine) *UserController {
	var controller = &UserController{Controller: core.CreateController(container, engine, "/user")}

	if err := container.Invoke(func(testService *UserService) {
		controller.testService = testService
	}); err != nil {
		log.Fatal(err)
	}

	controller.Router.GET("/create", controller.CreateUser)

	controller.Logger.Log("UserController initialized")

	return controller
}

func (r *UserController) CreateUser(c *gin.Context) {
	r.OkResponse(1).Reply(c)
}
