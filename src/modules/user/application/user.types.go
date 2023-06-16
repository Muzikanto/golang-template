package application

import (
	"context"
	"go-backend-template/src/modules/user/domain"
)

type UserService interface {
	Add(ctx context.Context, dto AddUserDto) (int64, error)
	Update(ctx context.Context, dto UpdateUserDto) error
	ChangePassword(ctx context.Context, dto ChangeUserPasswordDto) error
	GetById(ctx context.Context, userId int64) (domain.UserDto, error)
}
