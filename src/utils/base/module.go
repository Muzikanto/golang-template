package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"log"
	"reflect"
)

var moduleId int = 1

type IModule interface {
	Init(globalContainer *dig.Scope, engine *gin.Engine)
}

type Module struct {
	IModule
	engine      *gin.Engine
	controllers []IController
	modules     []IModule
	Providers   []interface{}
	Container   *dig.Scope
}

func (r *Module) Init(container *dig.Scope, engine *gin.Engine) {
	r.Container = container.Scope("module" + string(moduleId))
	moduleId++

	if r.engine == nil {
		r.engine = engine
	}

	var controllersCount = len(r.controllers)
	var modulesCount = len(r.modules)
	var providersCount = len(r.Providers)

	for i := 0; i < modulesCount; i++ {
		r.modules[i].Init(r.Container, r.engine)

		fmt.Print("[Module] - ", GetProviderKey(r.modules[i]), " initialized\n")
	}
	for i := 0; i < providersCount; i++ {
		var providerConstructor = r.Providers[i]
		var providerName = GetProviderKey(providerConstructor)

		// init
		var v = reflect.ValueOf(providerConstructor)

		if v.Kind() == reflect.Func {
			var r = v.Call([]reflect.Value{})
			var s = r[0]
			var t = reflect.TypeOf(s)

			if _, ok := t.MethodByName("Init"); ok {
				var initArgs = []reflect.Value{reflect.ValueOf(container)}
				s.MethodByName("Init").Call(initArgs)
			}
		} else {
			log.Fatalln(providerName + " is not a constructor")
		}

		// provide
		err := container.Provide(providerConstructor)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("[Provider] - ", providerName, " initialized\n")
	}
	for i := 0; i < controllersCount; i++ {
		r.controllers[i].Init(r.Container, r.engine)
		fmt.Print("[Controller] - ", GetProviderKey(r.controllers[i]), " initialized\n")
	}
}
func (r *Module) AddController(controller IController) {
	r.controllers = append(r.controllers, controller)
}

func (r *Module) AddProvider(value interface{}) {
	r.Providers = append(r.Providers, value)
}
