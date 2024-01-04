package main

import (
	"github.com/docker/docker/api/types"
	c8s "github.com/furon-kuina/cuternetes/pkg"
)

type Scenario []WorkerEvent

type WorkerEvent struct {
	container types.Container
	event     c8s.ContainerEvent
}
