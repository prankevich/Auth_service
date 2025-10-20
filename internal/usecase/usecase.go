package usecase

import (
	"github.com/prankevich/Auth_service/internal/adapter/driven/dbstore"
	"github.com/prankevich/Auth_service/internal/config"
	"github.com/prankevich/Auth_service/internal/port/usecase"
	authenticate "github.com/prankevich/Auth_service/internal/usecase/authenticator"
	usercreater "github.com/prankevich/Auth_service/internal/usecase/user_creator"
)

type UseCases struct {
	UserCreator   usecase.UserCreater
	Authenticator usecase.Authenticate
}

func New(cfg config.Config, store *dbstore.DBStore) *UseCases {
	return &UseCases{
		UserCreator:   usercreater.New(&cfg, store.UserStorage),
		Authenticator: authenticate.New(&cfg, store.UserStorage),
	}
}
