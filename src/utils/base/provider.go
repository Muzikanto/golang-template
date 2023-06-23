package base

import (
	"go.uber.org/dig"
)

type IProvider interface {
	Init(container *dig.Scope)
	GetName() string
}

type Provider struct {
	IProvider
	Name      string
	Container *dig.Scope
}

func (r *Provider) Init(container *dig.Scope) {
	r.Container = container
}

func (r *Provider) GetName() string {
	return r.Name
}
