package infrastructure

import (
	"context"
	"go-backend-template/src/modules/user/domain"
)

type UserRepository interface {
	Add(ctx context.Context, user domain.UserDomain) (int64, error)
	Update(ctx context.Context, user domain.UserDomain) (int64, error)
	GetById(ctx context.Context, userId int64) (domain.UserDomain, error)
	GetByEmail(ctx context.Context, email string) (domain.UserDomain, error)
}
