package user

import (
	"go-backend-template/src/utils/core"
	"go.uber.org/dig"
)

type UserService struct {
	*core.Service
}

func CreateTestService(container *dig.Scope) *UserService {
	var service = &UserService{Service: core.CreateService(container)}

	service.Logger.Log("UserService initialized")

	return service
}

func (r *UserService) CreateUser() int {
	return 1
}
