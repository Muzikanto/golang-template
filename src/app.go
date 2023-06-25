package src

import (
	"context"
	"go-backend-template/src/modules/config"
	"go-backend-template/src/modules/user"
	"go-backend-template/src/utils/db"
	"go-backend-template/src/utils/engine"
	"go-backend-template/src/utils/logger"
	"go.uber.org/dig"
	"log"
)

type App struct {
	Engine          *engine.Engine
	GlobalContainer *dig.Container
	Container       *dig.Scope
	Logger          *logger.Logger
	Config          *config.Config
	DB              *db.DbClient
}

func CreateApp() *App {
	var ctx = context.Background()

	var app = &App{}

	app.Engine = engine.CreateEngine()
	app.GlobalContainer = dig.New()
	app.Container = app.GlobalContainer.Scope("App")

	app.Logger = logger.CreateLogger("App", 0).
		Provide(app.Container)
	app.Config = config.
		CreateConfig(".env").
		Provide(app.Container)
	app.DB = db.CreateDB(ctx, app.Config.Database()).
		Connect().
		Provide(app.Container)

	return app
}

func (r *App) Init() {
	user.CreateModule(r.Container, r.Engine)
}

func (r *App) Start(addr string) {
	r.Logger.Log("Server listening at: " + addr)

	err := r.Engine.Run(addr)

	if err != nil {
		log.Fatal(err)
	}

}
