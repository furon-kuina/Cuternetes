package c8s

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
)

type ContainerSpec struct {
	Name  string
	Image string
	Cmd   strslice.StrSlice `json:",omitempty"`
}

type Spec struct {
	ContainerSpecs []ContainerSpec `yaml:"specs"`
}

type ContainerStatusResponse struct {
}

type ContainerState struct {
	Status ContainerStatus
}

type Worker struct {
	IsAvailable bool
	Containers  []types.Container
}

type C8sConfig struct {
	ApiServer ApiServerConfig `yaml:"api-server"`
	Workers   []WorkerConfig  `yaml:"workers"`
}

type WorkerConfig struct {
	Name string
	Url  string
}

type ApiServerConfig struct {
	Url string
}

type ContainerStatus string

const (
	Created    ContainerStatus = "created"
	Restarting ContainerStatus = "restarting"
	Running    ContainerStatus = "running"
	Removing   ContainerStatus = "removing"
	Paused     ContainerStatus = "paused"
	Exited     ContainerStatus = "exited"
	Dead       ContainerStatus = "dead"
)
