package dbstore

import (
	"database/sql"
	"errors"

	"github.com/prankevich/Auth_service/internal/errs"
)

func (u *UserStorage) translateError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotfound
	default:
		return err
	}
}
