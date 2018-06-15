package rancher

import (
	rancherapi "github.com/rancher/go-rancher/v2"
	"os"
	"fmt"
	"github.com/SantoDE/varaco/types"
)

var apiClient *rancherapi.RancherClient
var withoutPagination *rancherapi.ListOpts

func init() {
	rancherURL := os.Getenv("CATTLE_URL")
	accessKey := os.Getenv("CATTLE_ACCESS_KEY")
	secretKey := os.Getenv("CATTLE_SECRET_KEY")

	fmt.Printf("Injected Rancher Url %s \n", rancherURL)
	fmt.Printf("Injected Rancher Url %s \n", accessKey)
	fmt.Printf("Injected Rancher Url %s \n", secretKey)

	c, err := rancherapi.NewRancherClient(&rancherapi.ClientOpts{
		Url:       rancherURL,
		AccessKey: accessKey,
		SecretKey: secretKey,
	})

	if err != nil {
		fmt.Printf("Error Creating API Client: %s \n", err.Error())
	}

	apiClient = c;
}

func getSidekickContainer(selfContainerId string) (string, error) {

	fmt.Printf("Get Container By Uuid %s \n", selfContainerId)

	containers, err := apiClient.Container.List(withoutPagination)

	if err != nil {
		fmt.Printf("Cannot list containers for detecting Sidekick Container %s \n", err.Error())
		return "", err
	}

	self := filterContainerbyUuid(containers, selfContainerId)
	host := filterContainerbyDeploymentUnitUuid(containers, &self)

	return host.Uuid, nil
}

func filterContainerbyUuid(containers *rancherapi.ContainerCollection, uuid string) rancherapi.Container {

	for k := range containers.Data {
		container := containers.Data[k]

		if container.Uuid == uuid && container.State == "running"{
			fmt.Printf("Found a container by Uuid with id %s \n", container.Id)
			return container
		}
	}

	return *new(rancherapi.Container)
}

func filterContainerbyDeploymentUnitUuid(containers *rancherapi.ContainerCollection, self *rancherapi.Container) rancherapi.Container {

	for k := range containers.Data {
		container := containers.Data[k]

		if container.DeploymentUnitUuid == self.DeploymentUnitUuid && container.Id != self.Id && container.State == "running" {
			fmt.Printf("Found a container by Deployment Unit with id %s \n", container.Id)
			return container
		}
	}

	return *new(rancherapi.Container)
}

func executeCommand(containerId string) types.ExecuteCommand {

	containers, err := apiClient.Container.List(withoutPagination)

	target := filterContainerbyUuid(containers, containerId)

	fmt.Printf("Detected Target With id %s \n", target.Id)

	containerExec := new(rancherapi.ContainerExec)
	containerExec.AttachStdout = false
	containerExec.Command = []string{"varnish_reload_vcl"}

	cmd, err := apiClient.Container.ActionExecute(&target, containerExec)

	if err != nil {
		fmt.Printf("Cannot execute remote: %s \n", err.Error())
	}

	return types.ExecuteCommand{
		Url: cmd.Url,
		Token: cmd.Token,
	}
}