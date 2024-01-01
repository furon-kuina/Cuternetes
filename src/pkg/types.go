package c8s

import "github.com/docker/docker/api/types/strslice"

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

type ContainerStatus struct {
	IsAvailable bool
	Response    ContainerStatusResponse
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
