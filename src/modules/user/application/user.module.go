package application

import (
	userInfrastructure "go-backend-template/src/modules/user/infrastructure"
	"go-backend-template/src/utils/http"
)

type UserModule struct {
	UserController UserController
}

func NewUserModule(server *http.Server) *UserModule {
	userRepositoryOpts := userInfrastructure.UserRepositoryOpts{
		ConnManager: server.ConnManager,
	}
	userRepository := userInfrastructure.NewUserRepository(userRepositoryOpts)

	userServiceOpts := UserServiceOpts{
		TxManager:      server.TxManager,
		UserRepository: userRepository,
		Crypto:         server.Crypto,
	}
	userService := NewUserService(userServiceOpts)

	router := http.InitRouter(server, "/user")
	controller := UserController{
		Router:  router,
		Service: userService,
	}

	module := &UserModule{
		UserController: controller,
	}

	module.UserController.Init()

	print("User Module initialized\n")

	return module
}
