package providers

import (
	"context"
	"errors"
	"fmt"
	"github.com/virtual-kubelet/virtual-kubelet/internal/manager"
	"github.com/virtual-kubelet/virtual-kubelet/pkg/apis"
	"github.com/virtual-kubelet/virtual-kubelet/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
	"strings"
)

type TestProvider struct {
	client   		   apis.UnixSocketClient
	resourceManager    *manager.ResourceManager
	resourceGroup    	string
	nodeName         	string
	operatingSystem 	string
	clusterName        	string
	cpu               	string
	memory            	string
	pods               	string
	internalIP         	string
	daemonEndpointPort 	int32
	PodDataMap			map[string][]string
	PodLabelMap			map[string]map[string]string
	PodUIDMap			map[string]string
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
	if provider.PodDataMap == nil {
		provider.PodDataMap = make(map[string][]string)
	}
	if provider.PodUIDMap == nil{
		provider.PodUIDMap = make(map[string]string)
	}
	if provider.PodLabelMap == nil{
		provider.PodLabelMap = make(map[string]map[string]string)
	}
	podName := provider.getPodName(pod.Name, pod.Namespace)
	provider.PodDataMap[podName] = make([]string, 0)
	provider.PodUIDMap[podName] = string(pod.ObjectMeta.UID)
	provider.PodLabelMap[podName] = pod.Labels
	for _, container := range provider.getContainers(pod){
		containerID, err := provider.client.CreateDockerContainer(ctx, container.PodName, container.PodNamespace, container.ContainerName, container.Image)
		if err != nil {
			return err
		}
		provider.PodDataMap[podName] = append(provider.PodDataMap[podName], containerID)
	}
	return nil
}

func (provider *TestProvider) DeletePod(ctx context.Context, pod *v1.Pod) error {
	delete(provider.PodDataMap, provider.getPodName(pod.Name, pod.Namespace))
	for _, container := range provider.getContainers(pod){
		err := provider.client.DeleteDockerContainer(ctx, container.PodName, container.PodNamespace, container.ContainerName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (provider *TestProvider) UpdatePod(ctx context.Context, pod *v1.Pod) error {
	labels := pod.Labels
	podName := provider.getPodName(pod.Name, pod.Namespace)
	provider.PodLabelMap[podName] = pod.Labels
	if labels["state"] == "checkpoint" {
		containers, err := reserveContainersArray(provider.getContainers(pod))
		if err != nil {
			return err
		}
		for _, container := range containers {
			exit, err := strconv.ParseBool(labels["exit"])
			if err != nil {
				return err
			}
			err = provider.client.CheckpointDockerContainer(ctx, container.PodName, container.PodNamespace, container.ContainerName, labels["checkpoint-id"], exit)
			if err != nil {
				return err
			}
		}
	}else if labels["state"] == "restore" {
		for _, container := range provider.getContainers(pod){
			err := provider.client.RestoreDockerContainer(ctx, container.PodName, container.PodNamespace, container.ContainerName, labels["checkpoint-id"])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (provider *TestProvider) GetPod(ctx context.Context, namespace, name string) (*v1.Pod, error) {
	pods, err := provider.GetPods(ctx)
	if err != nil {
		return nil, err
	}
	for _, pod := range pods {
		if pod.Name == name && pod.Namespace == namespace {
			return pod, nil
		}
	}
	return nil, nil
}

func (provider *TestProvider) GetPodStatus(ctx context.Context, namespace, name string) (*v1.PodStatus, error) {
	pod, err := provider.GetPod(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	if pod == nil {
		return nil, nil
	}

	return &pod.Status, nil
}

func (provider *TestProvider) GetPods(ctx context.Context) ([]*v1.Pod, error) {
	results := make([]*v1.Pod, len(provider.PodDataMap))
	for key, containers := range provider.PodDataMap {
		namespace, name := provider.getPodMessage(key)
		k8sContainers, k8sContainersStatus, err := provider.convertContainerToK8SContainers(ctx , containers)
		if err != nil {
			return nil, err
		}
		results = append(results, &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Namespace: namespace,
				ClusterName: provider.clusterName,
				UID: types.UID(provider.PodUIDMap[key]),
			},
			Spec: v1.PodSpec{
				NodeName: provider.nodeName,
				Volumes: []v1.Volume{},
				Containers: k8sContainers,
			},
			Status: v1.PodStatus{
				Phase: "Running",
				Conditions: []v1.PodCondition{
					v1.PodCondition{
						Type:               v1.PodReady,
						Status:             v1.ConditionTrue,
						LastTransitionTime: metav1.Time{},
					}, v1.PodCondition{
						Type:               v1.PodInitialized,
						Status:             v1.ConditionTrue,
						LastTransitionTime: metav1.Time{},
					}, v1.PodCondition{
						Type:               v1.PodScheduled,
						Status:             v1.ConditionTrue,
						LastTransitionTime: metav1.Time{},
					},
				},
				Message: "",
				Reason: "",
				HostIP: "",
				PodIP: provider.internalIP,
				ContainerStatuses: k8sContainersStatus,
			},
		})
	}

	return nil, nil
}



func(provider *TestProvider) convertContainerToK8SContainers(context context.Context, containers []string) ([]v1.Container, []v1.ContainerStatus, error) {
	resultContainers := make([]v1.Container, len(containers))
	resultContainersStatus := make([]v1.ContainerStatus, len(containers))
	for _, containerID := range containers {
		dockerContainer, err := provider.client.GetContainerMessage(context, containerID)
		if err != nil {
			return nil, nil, err
		}
		resultContainers = append(resultContainers, v1.Container{
			Name: dockerContainer.ID,
			Image: dockerContainer.Image,
			Command: []string{dockerContainer.Command},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{},
			},
		})
		resultContainersStatus = append(resultContainersStatus, v1.ContainerStatus{
			Name: dockerContainer.ID,
			ContainerID: dockerContainer.ID,
			Image: dockerContainer.Image,
			ImageID: dockerContainer.ImageID,
			Ready: true,
		})
	}
	return resultContainers, resultContainersStatus, nil
}

//func containerStateToPodPhase(state string) v1.PodPhase {
//	return ""
//}

func reserveContainersArray(containers []ContainerData) ([]ContainerData, error) {
	length := len(containers)
	if length < 1 {
		return nil, errors.New("Invalid length.")
	}
	if length > 1 {
		for i := 0; i < length/2; i++ {
			temp := containers[length-1-i]
			containers[length-1-i] = containers[i]
			containers[i] = temp
		}
	}
	return containers, nil
}

func (provider *TestProvider) getPodName(name, namespace string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func (provider *TestProvider) getPodMessage(podName string) (string, string) {
	strArray := strings.Split(podName, "/")
	return strArray[0], strArray[1]
}

type ContainerData struct {
	PodName 		string
	PodNamespace 	string
	ContainerName 	string
	Image 			string
}

func (provider *TestProvider) getContainers(pod *v1.Pod) []ContainerData {
	containerDatas := make([]ContainerData, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		containerDatas = append(containerDatas, ContainerData{
			PodName: pod.Name,
			PodNamespace: pod.Namespace,
			ContainerName: container.Name,
			Image: container.Image,
		})
	}
	return containerDatas
}