package infrastructure

import (
	"context"
	database "go-backend-template/src/internal/base/database/impl"
	"go-backend-template/src/internal/base/errors"
	"go-backend-template/src/modules/user/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type UserRepositoryOpts struct {
	ConnManager database.ConnManager
}

func NewUserRepository(opts UserRepositoryOpts) UserRepository {
	return &userRepository{
		ConnManager: opts.ConnManager,
	}
}

type userRepository struct {
	database.ConnManager
}

func (r *userRepository) Add(ctx context.Context, model domain.UserDomain) (int64, error) {
	sql, _, err := database.QueryBuilder.
		Insert("users").
		Rows(database.Record{
			"firstname": model.FirstName,
			"lastname":  model.LastName,
			"email":     model.Email,
			"password":  model.Password,
		}).
		Returning("user_id").
		ToSQL()

	if err != nil {
		return 0, errors.Wrap(err, errors.DatabaseError, "syntax error")
	}

	row := r.Conn(ctx).QueryRow(ctx, sql)

	if err := row.Scan(&model.Id); err != nil {
		return 0, parseAddUserError(&model, err)
	}

	return model.Id, nil
}

func (r *userRepository) Update(ctx context.Context, model domain.UserDomain) (int64, error) {
	sql, _, err := database.QueryBuilder.
		Update("users").
		Set(database.Record{
			"firstname": model.FirstName,
			"lastname":  model.LastName,
			"email":     model.Email,
			"password":  model.Password,
		}).
		Where(database.Ex{"user_id": model.Id}).
		Returning("user_id").
		ToSQL()

	if err != nil {
		return 0, errors.Wrap(err, errors.DatabaseError, "syntax error")
	}

	row := r.Conn(ctx).QueryRow(ctx, sql)

	if err := row.Scan(&model.Id); err != nil {
		return 0, parseUpdateUserError(&model, err)
	}

	return model.Id, nil
}

func (r *userRepository) GetById(ctx context.Context, userId int64) (domain.UserDomain, error) {
	sql, _, err := database.QueryBuilder.
		Select(
			"firstname",
			"lastname",
			"email",
			"password",
		).
		From("users").
		Where(database.Ex{"user_id": userId}).
		ToSQL()

	if err != nil {
		return domain.UserDomain{}, errors.Wrap(err, errors.DatabaseError, "syntax error")
	}

	row := r.Conn(ctx).QueryRow(ctx, sql)

	model := domain.UserDomain{Id: userId}

	err = row.Scan(
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.Password,
	)
	if err != nil {
		return domain.UserDomain{}, parseGetUserByIdError(userId, err)
	}

	return model, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	sql, _, err := database.QueryBuilder.
		Select(
			"user_id",
			"firstname",
			"lastname",
			"password",
		).
		From("users").
		Where(database.Ex{"email": email}).
		ToSQL()

	if err != nil {
		return domain.UserDomain{}, errors.Wrap(err, errors.DatabaseError, "syntax error")
	}

	row := r.Conn(ctx).QueryRow(ctx, sql)

	model := domain.UserDomain{Email: email}

	err = row.Scan(
		&model.Id,
		&model.FirstName,
		&model.LastName,
		&model.Password,
	)
	if err != nil {
		return domain.UserDomain{}, parseGetUserByEmailError(email, err)
	}

	return model, nil
}

func parseAddUserError(user *domain.UserDomain, err error) error {
	pgError, isPgError := err.(*pgconn.PgError)

	if isPgError && pgError.Code == pgerrcode.UniqueViolation {
		switch pgError.ConstraintName {
		case "users_email_key":
			return errors.Wrapf(err, errors.AlreadyExistsError, "user with email \"%s\" already exists", user.Email)
		default:
			return errors.Wrapf(err, errors.DatabaseError, "add user failed")
		}
	}

	return errors.Wrapf(err, errors.DatabaseError, "add user failed")
}

func parseUpdateUserError(user *domain.UserDomain, err error) error {
	pgError, isPgError := err.(*pgconn.PgError)

	if isPgError && pgError.Code == pgerrcode.UniqueViolation {
		return errors.Wrapf(err, errors.AlreadyExistsError, "user with email \"%s\" already exists", user.Email)
	}

	return errors.Wrapf(err, errors.DatabaseError, "update user failed")
}

func parseGetUserByIdError(userId int64, err error) error {
	pgError, isPgError := err.(*pgconn.PgError)

	if isPgError && pgError.Code == pgerrcode.NoDataFound {
		return errors.Wrapf(err, errors.NotFoundError, "user with id \"%d\" not found", userId)
	}
	if err.Error() == "no rows in result set" {
		return errors.Wrapf(err, errors.NotFoundError, "user with id \"%d\" not found", userId)
	}

	return errors.Wrap(err, errors.DatabaseError, "get user by id failed")
}

func parseGetUserByEmailError(email string, err error) error {
	pgError, isPgError := err.(*pgconn.PgError)

	if isPgError && pgError.Code == pgerrcode.NoDataFound {
		return errors.Wrapf(err, errors.NotFoundError, "user with email \"%s\" not found", email)
	}
	if err.Error() == "no rows in result set" {
		return errors.Wrapf(err, errors.NotFoundError, "user with email \"%s\" not found", email)
	}

	return errors.Wrap(err, errors.DatabaseError, "get user by email failed")
}
