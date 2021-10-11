package providers

import (
	"context"
	"github.com/virtual-kubelet/virtual-kubelet/internal/manager"
	"github.com/virtual-kubelet/virtual-kubelet/pkg/config"
	v1 "k8s.io/api/core/v1"
)

type TestProvider struct {
	resourceManager    *manager.ResourceManager
	resourceGroup      string
	nodeName           string
	operatingSystem    string
	clusterName        string
	cpu                string
	memory             string
	pods               string
	internalIP         string
	daemonEndpointPort int32
}

func NewTestProvider(rm *manager.ResourceManager, nodeName, operatingSystem string, internalIP string, daemonEndpointPort int32) (*TestProvider, error) {
	provider := &TestProvider{}
	defaultConfig := config.GetDefaultConfig()

	provider.cpu = defaultConfig.CPU
	provider.memory = defaultConfig.Memory
	provider.clusterName = defaultConfig.ClusterName
	provider.pods = defaultConfig.Pods
	provider.resourceManager = rm
	provider.daemonEndpointPort = daemonEndpointPort
	provider.nodeName = nodeName
	provider.operatingSystem = operatingSystem
	provider.internalIP = internalIP

	return provider, nil
}

func (provider *TestProvider) CreatePod(ctx context.Context, pod *v1.Pod) error {

	return nil
}

func (provider *TestProvider) getContainers(pod *v1.Pod, init bool)   {

}