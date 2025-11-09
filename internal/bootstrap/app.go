package bootstrap

import (
	"context"
	"errors"
	_ "errors"
	"github.com/prankevich/Auth_service/internal/usecase"
	"github.com/prankevich/Auth_service/pkg/notification"

	"github.com/prankevich/Auth_service/internal/config"

	"net/http"
	"time"
)

const gracefulDeadline = 5 * time.Second

type App struct {
	cfg      config.Config
	rest     *http.Server
	useCases *usecase.UseCases
	teardown []func()
}

func New(cfg config.Config, producer notification.Producer) *App {
	app := initLayers(cfg)

	return app
}

func (app *App) Run(ctx context.Context) {
	go func() {
		if err := app.rest.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	for i := range app.teardown {
		app.teardown[i]()
	}
}

func (app *App) HTTPHandler() http.Handler {
	return app.rest.Handler
}
func (app *App) UseCases() *usecase.UseCases {
	return app.useCases
}
