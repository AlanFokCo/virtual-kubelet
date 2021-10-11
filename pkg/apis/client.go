package apis

import (
	tool "github.com/docker/docker/client"
)

type UnixSocketClient struct {
	Client *tool.Client
}

func NewUnixSocketClient() (*UnixSocketClient, error) {
	client, err := tool.NewClientWithOpts(tool.FromEnv)
	if err != nil {
		return nil, err
	}
	return &UnixSocketClient{
		Client: client,
	}, nil
}