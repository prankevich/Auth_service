package usecase

import (
	"auth_service/internal/adapter/driven/dbstore"
	"auth_service/internal/config"
	"auth_service/internal/port/usecase"
	authenticate "auth_service/internal/usecase/authenticator"
	usercreater "auth_service/internal/usecase/user_creator"
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
