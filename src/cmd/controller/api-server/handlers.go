package main

import (
	"encoding/json"
	"io"
	"net/http"

	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
)

func getContainersHandler(c echo.Context) error {
	containerStatuses := []*c8s.ContainerStatus{}
	for _, worker := range c8sConfig.Workers {
		resp, err := http.Get(worker.Url)
		status := new(c8s.ContainerStatus)
		if err != nil {
			status.IsAvailable = false
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			status.IsAvailable = false
			continue
		}
		err = json.Unmarshal(body, &status.Response)
		containerStatuses = append(containerStatuses, status)
	}
	return c.JSON(http.StatusOK, containerStatuses)
}

// putHandler accepts `kubectl apply`
func putHandler(c echo.Context) error {
	spec := new(c8s.Spec)
	if err := c.Bind(spec); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, spec)
}
