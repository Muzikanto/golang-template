package main

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/base"
	"go.uber.org/dig"
	"log"
)

// test module

type TestModule struct {
	*base.Module
}

func (r TestModule) Init(container *dig.Scope, engine *gin.Engine) {
	// providers
	r.Module.AddProvider(func() TestService {
		return TestService{Provider: &base.Provider{}}
	})

	// controllers
	var controller = TestController{Controller: &base.Controller{}}
	r.AddController(&controller)

	r.Module.Init(container, engine)
}

// test controller

type TestController struct {
	*base.Controller
	testService TestService
}

func (r *TestController) Init(container *dig.Scope, engine *gin.Engine) {
	r.Controller.Init(container, engine)
	r.Engine = engine

	if err := container.Invoke(func(testService TestService) {
		r.testService = testService

		if r.Engine != nil {
			r.Engine.GET("/test", r.CreateUser)
		} else {
			log.Fatal()
		}
	}); err != nil {
		log.Fatal(err)
	}
}

func (r *TestController) CreateUser(c *gin.Context) {
	c.JSON(200, nil)
}

// test provider

type TestService struct {
	*base.Provider
}

func (r *TestService) CreateUser() int {
	return 1
}

// main

func main() {
	var container = dig.New()

	var app = base.App{}

	app.Init()

	app.AddModule(TestModule{Module: &base.Module{}})
	app.InitModules(container)

	log.Fatal(app.Listen("127.0.0.1:3000"))
}
