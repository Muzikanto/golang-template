package user

import (
	"github.com/gin-gonic/gin"
	userDto "go-backend-template/src/modules/user/dto"
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
	//controller.Router.GET("/create", controller.Authorize, controller.AuthGuard, controller.RoleGuard([]string{"Admin"}), controller.CreateUser)

	controller.Logger.Log("UserController initialized")

	return controller
}

func (r *UserController) CreateUser(c *gin.Context) {
	var dto = userDto.CreateUserDto{}

	if err := r.Bind(&dto, c); err != nil {
		r.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	r.OkResponse(dto).Reply(c)
}
