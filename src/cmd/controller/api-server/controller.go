package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

func (c *ContainerController) Run(ctx context.Context, workerNum int) error {
	for {
		for _, worker := range c8sConfig.Workers {
			resp, err := http.Get(worker.Url + "/containers")
			if err != nil {
				fmt.Printf("got error on request: %v", err)
			}
			fmt.Printf("response: %#v", resp)
			defer resp.Body.Close()
			payload, err := io.ReadAll(resp.Body)
			fmt.Println(payload)
			if err != nil {
				fmt.Printf("Failed to read response body: %v", err)
			}
			fmt.Println(payload)
		}
	}
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
