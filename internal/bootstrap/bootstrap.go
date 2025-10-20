package bootstrap

import (
	"auth_service/internal/adapter/driven/dbstore"
	"auth_service/internal/config"
	"auth_service/internal/usecase"
	"context"
	"fmt"
)

func initLayers(cfg config.Config) *App {
	teardown := make([]func(), 0)
	db, err := initDB(*cfg.Postgres)
	if err != nil {
		panic(err)
	}

	storage := dbstore.New(db)
	teardown = append(teardown, func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
			//log.Error(err.Error())
		}
	})
	uc := usecase.New(cfg, storage)

	httpSrv := initHTTPService(&cfg, uc)

	teardown = append(teardown,
		func() {
			ctxShutDown, cancel := context.WithTimeout(context.Background(), gracefulDeadline)
			defer cancel()
			if err := httpSrv.Shutdown(ctxShutDown); err != nil {
				//log.Error(fmt.Sprintf("server Shutdown Failed:%s", err))
				return
			}
		},
	)

	return &App{
		cfg:      cfg,
		rest:     httpSrv,
		teardown: teardown,
	}
}
