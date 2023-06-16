package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/crypto"
	"go-backend-template/src/utils/database"
	databaseImpl "go-backend-template/src/utils/database/impl"
)

type Config interface {
	DetailedError() bool
	Address() string
}

type ServerOpts struct {
	ConnManager databaseImpl.ConnManager
	Crypto      crypto.Crypto
	Config      Config
}

func NewServer(opts ServerOpts) *Server {
	gin.SetMode(gin.ReleaseMode)

	server := &Server{
		Engine:      gin.New(),
		config:      opts.Config,
		Crypto:      opts.Crypto,
		ConnManager: opts.ConnManager,
	}

	return server
}

type Server struct {
	Engine      *gin.Engine
	config      Config
	Crypto      crypto.Crypto
	ConnManager databaseImpl.ConnManager
	TxManager   database.TxManager
}

func (s Server) Listen() error {
	fmt.Printf("API server listening at: %s\n\n", s.config.Address())
	return s.Engine.Run(s.config.Address())
}
