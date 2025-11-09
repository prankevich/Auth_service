package authenticate

import (
	"context"
	"errors"

	"github.com/prankevich/Auth_service/internal/config"
	"github.com/prankevich/Auth_service/internal/domain"
	"github.com/prankevich/Auth_service/internal/errs"
	"github.com/prankevich/Auth_service/internal/port/driven"
	"github.com/prankevich/Auth_service/utils"
)

type UseCase struct {
	cfg         *config.Config
	userStorage driven.UserStorage
}

func New(cfg *config.Config, userStorage driven.UserStorage) *UseCase {
	return &UseCase{
		cfg:         cfg,
		userStorage: userStorage,
	}
}

func (u *UseCase) Authenticate(ctx context.Context, user domain.User) (int, domain.Role, error) {
	userFromDB, err := u.userStorage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotfound) {
			return 0, "", errs.ErrUserNotFound
		}

		return 0, "", err
	}
	user.Password, err = utils.GenerateHash(user.Password)
	if err != nil {
		return 0, "", err
	}

	// проверить правильно ли он указал пароль
	if userFromDB.Password != user.Password {
		return 0, "", errs.ErrIncorrectUsernameOrPassword
	}

	return userFromDB.ID, userFromDB.Role, nil
}
