package apis

import (
	tool "github.com/docker/go-docker"
	"net/http"
)

type UnixSocketClient struct {
	Client *tool.Client
}

func NewUnixSocketClient(host string, version string, client *http.Client, httpHeaders map[string]string) (*UnixSocketClient, error) {
	dockerClient, err := tool.NewClient(host, version, client, httpHeaders)
	if err != nil {
		return nil, err
	}
	return &UnixSocketClient{
		Client: dockerClient,
	}, nil
}