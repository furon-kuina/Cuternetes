package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/docker/docker/client"
	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
	"gopkg.in/yaml.v3"
)

var (
	workerName string
	cli        *client.Client
	c8sConfig  c8s.C8sConfig
	port       string
)

const c8sConfigPath string = "../config.yaml"

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	f, err := os.Open(c8sConfigPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	data, err := io.ReadAll(f)
	yaml.Unmarshal(data, &c8sConfig)
	workerName = os.Getenv("WORKER_NAME")
	for _, worker := range c8sConfig.Workers {
		log.Printf(worker.Url)
		log.Printf(worker.Name)
		log.Printf(workerName)
		if worker.Name == workerName {
			tmp := strings.Split(worker.Url, ":")
			port = fmt.Sprint(tmp[len(tmp)-1])
		}
	}
}

func setRoutes(e *echo.Echo) {
	e.GET("/", getDefaultHandler)
	e.GET("/containers", getContainersHandler)
	e.POST("/create", postCreateHandler)
	e.POST("/delete", postDeleteHandler)
}

func main() {
	e := echo.New()
	setRoutes(e)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	e.Use(slogecho.New(logger))
	log.Printf("Hello, world!")
	e.Logger.Fatal(e.Start(":" + fmt.Sprint(port)))
}

// TODO: keep sending heartbeat in case of network error
