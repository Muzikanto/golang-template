package main

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"reflect"
)

type ICat interface {
	Meow(name string)
}

type Cat struct {
	ICat
	value int
}

func (c *Cat) Meow(name string) {
	fmt.Println("Meow", name)
}

func main() {
	var container = dig.New()

	AddCat(container, func() *Cat {
		var cat = &Cat{value: 2}
		return cat
	})

	err := container.Invoke(func(cat *Cat) {
		cat.Meow("Cat")
	})
	if err != nil {
		log.Fatal()
	}
}

func AddCat(container *dig.Container, constructor interface{}) {
	err := container.Provide(constructor)

	if err != nil {
		log.Fatal(err)
	}

	var t = reflect.TypeOf(constructor)
	var v = reflect.ValueOf(constructor)

	if v.Kind() == reflect.Func {
		var r = v.Call([]reflect.Value{})
		var s = r[0]
		var args = []reflect.Value{reflect.ValueOf("test")}
		var m = s.MethodByName("Test")

		if _, ok := t.MethodByName("test"); ok {
			m.Call(args)
		}
	}
}
