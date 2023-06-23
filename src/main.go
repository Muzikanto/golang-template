package main

import (
	"context"
	"go-backend-template/src/config"
	"go-backend-template/src/modules/test"
	userApplication "go-backend-template/src/modules/user/application"
	"go-backend-template/src/utils/crypto/impl"
	"go-backend-template/src/utils/database/impl"
	"go-backend-template/src/utils/http"
	"go.uber.org/dig"
	"log"
)

func main() {
	container := dig.New()
	ctx := context.Background()

	// config
	conf := config.ParseEnv(container, ".env")

	// db
	dbClient := database.NewClient(ctx, conf.Database())
	if err := container.Provide(func() *database.Client {
		return dbClient
	}); err != nil {
		log.Fatal(err)
		return
	}

	err := dbClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer dbClient.Close()

	// crypto
	crypto := impl.NewCrypto()
	container.Provide(crypto)

	dbService := database.NewService(dbClient)

	serverOpts := http.ServerOpts{
		Crypto:      crypto,
		Config:      conf.HTTP(),
		ConnManager: dbService,
	}
	server := http.NewServer(serverOpts)

	//
	userApplication.NewUserModule(
		server,
	)

	test.CreateController(server.Engine)

	log.Fatal(server.Listen())
}
