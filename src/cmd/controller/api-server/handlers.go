package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/docker/docker/api/types"
	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
)

func getContainersHandler(c echo.Context) error {
	workers := make([]c8s.Worker, len(c8sConfig.Workers))
	for i, worker := range c8sConfig.Workers {
		resp, err := http.Get(worker.Url + "/containers")
		if err != nil {
			workers[i].IsAvailable = false
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			workers[i].IsAvailable = false
			continue
		}
		var containers []types.Container
		err = json.Unmarshal(body, &containers)
		if err != nil {
			workers[i].IsAvailable = false
		}
		workers[i].Containers = containers
	}
	return c.JSON(http.StatusOK, workers)
}

// putHandler accepts `kubectl apply`
func putHandler(c echo.Context) error {
	spec := new(c8s.Spec)
	if err := c.Bind(spec); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, spec)
}

// receives events from workers
func postEventHandler(c echo.Context) error {
	return nil
}
