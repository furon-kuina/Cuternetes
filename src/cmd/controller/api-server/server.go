package main

import (
	"context"
	"io"
	"log"
	"os"

	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

var (
	c8sConfig c8s.C8sConfig
)

const c8sConfigPath string = "../../config.yaml"

func init() {
	f, err := os.Open(c8sConfigPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	data, err := io.ReadAll(f)
	yaml.Unmarshal(data, &c8sConfig)
}

func main() {
	// e := echo.New()
	// setRoutes(e)
	// e.Use(middleware.Logger())
	// e.Logger.Fatal(e.Start(":" + string(c8sConfig.ApiServer.Port)))

	ctx := context.Background()
	c := NewContainerController(ctx)
	c.Run(ctx)
}

func setRoutes(e *echo.Echo) {
	e.GET("/containers", getContainersHandler)
	e.PUT("/", putHandler)
}
