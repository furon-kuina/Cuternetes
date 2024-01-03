package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
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

func (c *ContainerController) Run(ctx context.Context) error {
	for {
		workerContainers := c.GetContainerStatus(ctx)
		for k, v := range workerContainers {
			fmt.Printf("worker: %s, containers: %+v", k, v)
		}
		time.Sleep(5 * time.Second)
	}
}

// never returns error
// if something went wrong with communication,
// it just fills the corresponding slot with nil
func (c *ContainerController) GetContainerStatus(ctx context.Context) map[string][]types.Container {
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

// func (c *ContainerController) Reconcile(ctx context.Context) error {

// }
