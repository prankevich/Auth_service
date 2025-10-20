package bootstrap

import (
	"context"
	"fmt"
	"github.com/prankevich/Auth_service/internal/adapter/driven/dbstore"

	"github.com/prankevich/Auth_service/internal/config"
	"github.com/prankevich/Auth_service/internal/usecase"
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
