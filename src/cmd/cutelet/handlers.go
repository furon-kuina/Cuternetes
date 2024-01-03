package main

import (
	"context"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
)

func postCreateHandler(c echo.Context) error {
	container := new(c8s.ContainerSpec)
	if err := c.Bind(container); err != nil {
		return err
	}

	if err := createContainer(container); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, container)
}

func getContainersHandler(c echo.Context) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	log.Println(containers)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, containers)
}

func getDefaultHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello, cutelet!")
}
