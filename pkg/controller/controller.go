package controller

import (
	"Stale-purger/pkg/config"
	"context"

	"github.com/sirupsen/logrus"
)

type Component interface {
	Start(ctx context.Context)
	Name() string
}

type Controller struct {
	Components []Component
	Config     config.Config
	Logger     *logrus.Entry
}

func NewController(config config.Config, logger *logrus.Entry) *Controller {
	return &Controller{
		Logger: logger,
	}
}

func (c *Controller) AddComponent(component Component) {
	c.Components = append(c.Components, component)
}

func (c *Controller) Start(ctx context.Context) {
	c.Logger.Info("Starting components...")
	for _, comp := range c.Components {
		go func(comp Component) {
			c.Logger.Printf("Starting %s...", comp.Name())
			comp.Start(ctx)
		}(comp)
	}
}
