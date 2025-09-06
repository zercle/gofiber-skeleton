package container

import (
	"context"
	"log"

	"go.uber.org/fx"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
)

type Container struct {
	app *fx.App
}

func NewContainer() *Container {
	app := fx.New(
		fx.Provide(
			config.Load,
			database.NewDatabase,
		),
		fx.Invoke(func(cfg *config.Config, db *database.Database) {
			log.Printf("Application initialized with config: %s:%s", cfg.App.Name, cfg.App.Port)
		}),
	)

	return &Container{app: app}
}

func (c *Container) Start(ctx context.Context) error {
	return c.app.Start(ctx)
}

func (c *Container) Stop(ctx context.Context) error {
	return c.app.Stop(ctx)
}

func (c *Container) App() *fx.App {
	return c.app
}