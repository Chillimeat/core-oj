package dockerx

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	client "github.com/docker/docker/client"
)

// Client is the alias of docker client
type Client = client.Client

// Connect to Remote Docker Server
func Connect(host, version string) (*Client, error) {
	// ctx := context.Background()

	cli, err := client.NewClient(host, version, nil, nil)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// ListAllContainers Test Docker client works well
func ListAllContainers(cli *Client) ([]types.Container, error) {
	return cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
}

func SearchContainerByName(cli *Client, name string) (*types.Container, error) {
	var filter = filters.NewArgs()
	filter.Add("name", "^"+name+"$")

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return nil, nil
	}
	return &containers[0], nil
}

func SearchContainerByID(cli *Client, id string) (*types.Container, error) {
	var filter = filters.NewArgs()
	filter.Add("id", "^"+id+"$")

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return nil, nil
	}
	return &containers[0], nil
}
