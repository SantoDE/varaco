package rancher

import (
	"github.com/SantoDE/varaco/configuration"
	"github.com/rancher/go-rancher-metadata/metadata"
	"fmt"
	"github.com/SantoDE/varaco/types"
	"context"
)

const (
	RefreshSeconds = 60
)

var metadataClient metadata.Client

type Rancher struct {
	Config *configuration.Rancher
}

func init() {
	metadataServiceURL := fmt.Sprintf("http://rancher-metadata.rancher.internal/latest")
	mc, err := metadata.NewClientAndWait(metadataServiceURL)

	if err != nil {
		fmt.Printf("Error Creating Metadata Client %s \n", err.Error())
	}

	metadataClient = mc
}

func listContainerIps(serviceToWatch string) []types.ContainerData {

	var containersData []types.ContainerData

	stack, err := metadataClient.GetSelfStack()

	services := stack.Services

	for k := range services {
		service := services[k]

		if service.Name == serviceToWatch {
			containers := service.Containers

			for sk := range containers {
				container := containers[sk];

				containerData := new(types.ContainerData)
				containerData.Ip = container.PrimaryIp
				containerData.Id = container.UUID
				containersData = append(containersData, *containerData)
			}
		}
	}

	if err != nil {
		fmt.Printf("Error %s \n", err.Error())
	}

	return containersData
}

func getSelfContainerId() (string, error) {
	self, err := metadataClient.GetSelfContainer()

	if err != nil {
		return "", err
	}

	return self.UUID, nil
}

func longPoll(updateConfiguration func(string)) {

	fmt.Printf("In longPoll \n")
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fmt.Printf("In gofunc() on change \n")
		metadataClient.OnChange(RefreshSeconds, updateConfiguration)
	}()
}