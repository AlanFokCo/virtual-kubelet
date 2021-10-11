package apis

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (client *UnixSocketClient) CreateDockerContainer(context context.Context, name, namespace, containerName, image string) error {
	if DockerMap == nil {
		DockerMap = make(map[string]string)
	}

	resp, err := client.Client.ContainerCreate(context, &container.Config{
		Image: image,
		Cmd: []string{},
	}, nil, nil, GetContainerName(name, namespace, containerName))
	if err != nil {
		return err
	}

	DockerMap[GetContainerName(name, namespace, containerName)] = resp.ID

	return nil
}

func (client *UnixSocketClient) DeleteDockerContainer(context context.Context, name, namespace, containerName string) error {
	err := client.Client.ContainerRemove(context, DockerMap[GetContainerName(name, namespace, containerName)], types.ContainerRemoveOptions{})
	if err != nil{
		return err
	}
	return nil
}

func (client *UnixSocketClient) CheckpointDockerContainer(context context.Context, name, namespace, containerName, checkpointID string, exit bool) error {
	err := client.Client.CheckpointCreate(context, DockerMap[GetContainerName(name, namespace, containerName)], types.CheckpointCreateOptions{
		CheckpointID: checkpointID,
		CheckpointDir: "/opt/checkpoint/" + checkpointID,
		Exit: exit,
	})
	if err != nil{
		return err
	}
	return nil
}

func (client *UnixSocketClient) RestoreDockerContainer(context context.Context, name, namespace, containerName, checkpointID string) error {
	err := client.Client.ContainerStart(context, DockerMap[GetContainerName(name, namespace, containerName)], types.ContainerStartOptions{
		CheckpointID: checkpointID,
		CheckpointDir: "/opt/checkpoint/" + checkpointID,
	})
	if err != nil{
		return err
	}
	return nil
}

func GetContainerName(name, namespace, containerName string) string {
	return fmt.Sprintf("%s-%s-%s", namespace, name, containerName)
}