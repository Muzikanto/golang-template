package application

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/internal/base/crypto"
	"go-backend-template/src/internal/base/database"
	databaseImpl "go-backend-template/src/internal/base/database/impl"
	userInfrastructure "go-backend-template/src/modules/user/infrastructure"
)

type UserModule struct {
	UserController UserController
}

func NewUserModule(engine *gin.Engine, connManager databaseImpl.ConnManager, txManager database.TxManager, crypto crypto.Crypto) *UserModule {
	userRepositoryOpts := userInfrastructure.UserRepositoryOpts{
		ConnManager: connManager,
	}
	userRepository := userInfrastructure.NewUserRepository(userRepositoryOpts)

	userServiceOpts := UserServiceOpts{
		TxManager:      txManager,
		UserRepository: userRepository,
		Crypto:         crypto,
	}
	userService := NewUserService(userServiceOpts)

	module := &UserModule{
		UserController: UserController{
			Engine:  engine,
			Service: userService,
		},
	}

	module.UserController.Init()

	return module
}
