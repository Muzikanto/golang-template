package main

import (
	"context"
	"go-backend-template/src/api/cli"
	"go-backend-template/src/api/http"
	"go-backend-template/src/internal/base/crypto/impl"
	"go-backend-template/src/internal/base/database/impl"
	userApplication "go-backend-template/src/modules/user/application"
	"log"
)

func main() {
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
	dbService := database.NewService(dbClient)

	serverOpts := http.ServerOpts{
		Crypto:      crypto,
		Config:      conf.HTTP(),
		ConnManager: dbService,
	}
	server := http.NewServer(serverOpts)

	//
	userApplication.NewUserModule(
		server.Engine,
		dbService,
		dbService,
		crypto,
	)

	log.Fatal(server.Listen())
}
