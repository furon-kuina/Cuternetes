package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	c8s "github.com/furon-kuina/cuternetes/pkg"
)

func createContainer(c *c8s.ContainerSpec) (err error) {
	defer c8s.Wrap(&err, "createContainer(%q)", c)
	ctx := context.Background()
	err = pullImage(ctx, c.Image)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: c.Image,
			Cmd:   c.Cmd,
		}, nil, nil, nil, c.Name)

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	log.Printf("Started container with ID %s", resp.ID)
	return nil
}

func deleteContainer(containerId string) (err error) {
	defer c8s.Wrap(&err, "deleteContainer(%q)", containerId)
	ctx := context.Background()
	if err := cli.ContainerStop(ctx, containerId, container.StopOptions{}); err != nil {
		return err
	}
	if err := cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	log.Printf("Deleted container with ID %s", containerId)
	return nil
}

func pullImage(ctx context.Context, image string) (err error) {
	defer c8s.Wrap(&err, "pullImage(%q)", image)
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("image pull failed: %v", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}
