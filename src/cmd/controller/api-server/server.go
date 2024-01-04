package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

var (
	c8sConfig     c8s.C8sConfig
	spec          c8s.Spec
	workerConfigs = make(map[string]c8s.WorkerConfig)
)

const c8sConfigPath string = "/home/tenma/src/github.com/furon-kuina/Cuternetes/src/cmd/config.yaml"
const c8sSpecPath string = "/home/tenma/src/github.com/furon-kuina/Cuternetes/src/cmd/cutectl/specs.yaml"

func init() {
	f, err := os.Open(c8sConfigPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	data, err := io.ReadAll(f)
	yaml.Unmarshal(data, &c8sConfig)
	fmt.Printf("config: %+v\n", c8sConfig)

	for _, workerConfig := range c8sConfig.Workers {
		workerConfigs[workerConfig.Name] = workerConfig
	}

	f, err = os.Open(c8sSpecPath)
	if err != nil {
		log.Fatalf("failed to open spec file: %v", err)
	}
	data, err = io.ReadAll(f)
	yaml.Unmarshal(data, &spec)
	fmt.Printf("specs: %+v\n", spec)
}

func validateSpecs(specs c8s.Spec) bool {
	// names need to be unique
	return true
}

func main() {
	// e := echo.New()
	// setRoutes(e)
	// e.Use(middleware.Logger())
	// e.Logger.Fatal(e.Start(":" + string(c8sConfig.ApiServer.Port)))

	ctx := context.Background()
	c := NewContainerController(ctx)
	c.Reconcile()
}

func setRoutes(e *echo.Echo) {
	e.GET("/containers", getContainersHandler)
	e.PUT("/", putHandler)
	e.POST("/events", postEventHandler)
}
