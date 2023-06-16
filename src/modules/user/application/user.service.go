package application

import (
	"context"
	"go-backend-template/src/modules/user/domain"
	"go-backend-template/src/modules/user/infrastructure"
	"go-backend-template/src/utils/crypto"
	"go-backend-template/src/utils/database"
)

type UserServiceOpts struct {
	TxManager      database.TxManager
	UserRepository infrastructure.UserRepository
	Crypto         crypto.Crypto
}

func NewUserService(opts UserServiceOpts) UserService {
	return &userService{
		TxManager:      opts.TxManager,
		UserRepository: opts.UserRepository,
		Crypto:         opts.Crypto,
	}
}

type userService struct {
	database.TxManager
	infrastructure.UserRepository
	crypto.Crypto
}

func (u *userService) Add(ctx context.Context, in AddUserDto) (userId int64, err error) {
	model, err := in.ToDomain()
	if err != nil {
		return 0, err
	}
	if err := model.HashPassword(u.Crypto); err != nil {
		return 0, err
	}

	// Transaction demonstration
	err = u.RunTx(ctx, func(ctx context.Context) error {
		userId, err = u.UserRepository.Add(ctx, model)
		if err != nil {
			return err
		}
		model.Id = userId

		userId, err = u.UserRepository.Update(ctx, model)
		if err != nil {
			return err
		}
		return nil
	})

	return userId, err
}

func (u *userService) Update(ctx context.Context, in UpdateUserDto) (err error) {
	model, err := u.UserRepository.GetById(ctx, in.Id)
	if err != nil {
		return err
	}
	err = model.Update(in.FirstName, in.LastName, in.Email)
	if err != nil {
		return err
	}
	_, err = u.UserRepository.Update(ctx, model)

	return err
}

func (u *userService) ChangePassword(ctx context.Context, in ChangeUserPasswordDto) (err error) {
	user, err := u.UserRepository.GetById(ctx, in.Id)
	if err != nil {
		return err
	}
	if err = user.ChangePassword(in.Password, u.Crypto); err != nil {
		return err
	}
	_, err = u.UserRepository.Update(ctx, user)

	return err
}

func (u *userService) GetById(ctx context.Context, userId int64) (out domain.UserDto, err error) {
	model, err := u.UserRepository.GetById(ctx, userId)
	if err != nil {
		return out, err
	}

	return out.ToResponse(model), nil
}
