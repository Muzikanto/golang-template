package main

import (
	"context"
	userApplication "go-backend-template/src/modules/user/application"
	"go-backend-template/src/utils/cli"
	"go-backend-template/src/utils/crypto/impl"
	"go-backend-template/src/utils/database/impl"
	"go-backend-template/src/utils/http"
	"go.uber.org/dig"
	"log"
)

func main() {
	container := dig.New()

	ctx := context.Background()
	parser := cli.NewParser()

	conf, err := parser.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbClient := database.NewClient(ctx, conf.Database())

	err = dbClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer dbClient.Close()

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

	log.Fatal(server.Listen())
}
