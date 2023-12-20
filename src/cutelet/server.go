package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("Failed to create a Docker client: %v", err)
		return
	}

	ctx := context.Background()

	imageName := "alpine"
	err = pullImage(ctx, cli, imageName)
	if err != nil {
		log.Fatalf("Failed to pull image: %v", err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"}}, nil, nil, nil, "")

	if err != nil {
		log.Fatalf("Failed to create a container: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to start a container: %v", err)
		return
	}

	fmt.Printf("Started container with ID %s", resp.ID)
}

func pullImage(ctx context.Context, cli *client.Client, image string) error {
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("Failed to pull error: %v", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}
