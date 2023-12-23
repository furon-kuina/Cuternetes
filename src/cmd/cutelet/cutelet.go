package main

import (
	"log"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var cli *client.Client

func configure() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}
}

func setRoutes(e *echo.Echo) {
	e.POST("/create", postCreateHandler)
}

func main() {
	configure()
	e := echo.New()
	setRoutes(e)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":3333"))
}
