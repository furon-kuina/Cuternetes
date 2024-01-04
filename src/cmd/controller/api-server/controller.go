package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	c8s "github.com/furon-kuina/cuternetes/pkg"
)

type ContainerController struct {
	ctx context.Context
	cli http.Client
}

func NewContainerController(ctx context.Context) *ContainerController {
	return &ContainerController{
		ctx: ctx,
	}
}

type Controller interface {
	Run(ctx context.Context, workerNum int) error
	Reconcile() error
}

func (c *ContainerController) Run() error {
	for {
		if err := c.Reconcile(); err != nil {
			return fmt.Errorf("reconcile failed: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (c *ContainerController) Reconcile() error {
	desired, err := c.FetchContainerSpecs()
	if err != nil {
		fmt.Printf("failed to fetch spec: %v", err)
	}
	fmt.Printf("desired state: %+v\n", desired)
	lacking := make([]c8s.ContainerSpec, 0)
	unnecessaryContainers := make([]types.Container, 0)
	workers := make([]string, 0)
	containersInWorker := c.GetContainerStatus()
	existing := make(map[string]string)
	for workerName, containers := range containersInWorker {
		for _, container := range containers {
			if _, ok := desired[container.Names[0]]; !ok {
				unnecessaryContainers = append(unnecessaryContainers, container)
				workers = append(workers, workerName)
			} else {
				existing[container.Names[0]] = workerName
			}
		}
	}
	for name, config := range desired {
		if _, ok := existing[name]; !ok {
			lacking = append(lacking, config)
		}
	}
	fmt.Printf("unnecessary containers: %+v\n", unnecessaryContainers)
	fmt.Printf("lacking containers: %+v\n", lacking)
	for i, container := range unnecessaryContainers {
		if err := c.DeleteContainer(container, workerConfigs[workers[i]]); err != nil {
			// TODO: handle error
		}
	}
	for _, container := range lacking {
		targetWorkerName := c.ScheduleContainer(container, containersInWorker)
		fmt.Printf("target worker is %s with config %+v\n", targetWorkerName, workerConfigs[targetWorkerName])
		if err := c.CreateContainerAt(workerConfigs[targetWorkerName], container); err != nil {
			return fmt.Errorf("failed to create contianer: %v", err)
		}
	}
	return nil
}

func (c *ContainerController) DeleteContainer(container types.Container, workerConfig c8s.WorkerConfig) error {
	fmt.Printf("Deleting container: %+v\n\n", container)
	data, _ := json.Marshal(container)
	resp, err := http.Post(workerConfig.Url+"/delete", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to delete container %+v in worker %s", container, workerConfig.Name)
	}
	fmt.Printf("response from worker %s: %+v\n", workerConfig.Name, resp)
	return nil
}

func (c *ContainerController) CreateContainerAt(workerConfig c8s.WorkerConfig, container c8s.ContainerSpec) error {
	data, err := json.Marshal(container)
	if err != nil {
		return fmt.Errorf("failed to marshal container config: %v\n", err)
	}
	fmt.Printf("sending create request to worker %s: %s\n", workerConfig.Name, string(data))
	resp, err := http.Post(workerConfig.Url+"/create", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}
	fmt.Printf("response from worker %s: %+v\n", workerConfig.Name, resp)
	return nil
}

func (c *ContainerController) ScheduleContainer(container c8s.ContainerSpec, containersInWorker map[string][]types.Container) string {
	minimum := 10000000
	argmin := ""
	for workerName, containers := range containersInWorker {
		fmt.Printf("worker %s has %d containers\n", workerName, len(containers))
		if len(containers) < minimum {
			argmin = workerName
		}
	}
	fmt.Printf("next container will run in worker %s\n", argmin)
	return argmin
}

func (c *ContainerController) FetchContainerSpecs() (map[string]c8s.ContainerSpec, error) {
	s := make(map[string]c8s.ContainerSpec)
	for _, specs := range spec.ContainerSpecs {
		s[specs.Name] = specs
	}
	return s, nil
}

// if something went wrong with communication,
// it just fills the corresponding slot with nil
func (c *ContainerController) GetContainerStatus() map[string][]types.Container {
	workerContainers := make(map[string][]types.Container)
	for _, worker := range c8sConfig.Workers {
		resp, err := http.Get(worker.Url + "/containers")
		if err != nil {
			workerContainers[worker.Name] = nil
			continue
		}
		defer resp.Body.Close()
		payload, err := io.ReadAll(resp.Body)
		if err != nil {
			workerContainers[worker.Name] = nil
			continue
		}
		var containers []types.Container
		json.Unmarshal(payload, &containers)
		fmt.Printf("containers of %s: %+v\n", worker.Name, containers)
		workerContainers[worker.Name] = containers
	}
	return workerContainers
}

// func (c *ContainerController) runWorker(ctx context.Context) {
// 	for c.processNextItem(ctx) {

// 	}
// }

// func (c *ContainerController) processNextItem(ctx context.Context) bool {
// 	obj, shutdown := c.workqueue.Get()

// 	if err := c.Reconcile(); err != nil {
// 		c.workqueue.Add()
// 	}
// }

// type WorkQueue struct {
// }

// func (wq *WorkQueue) Get() (obj, err) {

// }
