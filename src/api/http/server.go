package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-backend-template/src/internal/base/crypto"
	databaseImpl "go-backend-template/src/internal/base/database/impl"
	userApplication "go-backend-template/src/modules/user/application"
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
		crypto:      opts.Crypto,
		ConnManager: opts.ConnManager,
	}

	initRouter(server)

	return server
}

type Server struct {
	Engine      *gin.Engine
	config      Config
	crypto      crypto.Crypto
	UserService userApplication.UserService
	ConnManager databaseImpl.ConnManager
}

func (s Server) Listen() error {
	fmt.Printf("API server listening at: %s\n\n", s.config.Address())
	return s.Engine.Run(s.config.Address())
}
