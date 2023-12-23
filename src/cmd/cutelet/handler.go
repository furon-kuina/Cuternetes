package main

import (
	"net/http"

	c8s "github.com/furon-kuina/cuternetes/pkg"
	"github.com/labstack/echo/v4"
)

func postCreateHandler(c echo.Context) error {
	container := new(c8s.Container)
	if err := c.Bind(container); err != nil {
		return err
	}

	if err := createContainer(container); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, container)
}
