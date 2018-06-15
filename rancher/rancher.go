package rancher

import (
	"github.com/SantoDE/varaco/types"
	"fmt"
	"os"
)

func (r *Rancher) Provide(serviceToWatch string, varnishConfig chan<- types.VarnishConfiguration) {

	fmt.Printf("Provide Data.... \n")

	containers := listContainerIps(serviceToWatch)

	update := func(version string) {
		containers = listContainerIps(serviceToWatch)
		selfContainerId, err := getSelfContainerId()

		if err != nil {
			os.Exit(15)
			fmt.Printf("Error Getting Self Container Id %s \n", err.Error())
		}

		varnishHost, err := getSidekickContainer(selfContainerId)

		if err != nil {
			fmt.Printf("Error Getting Sidekick Container Ip %s \n", err.Error())
			os.Exit(15);
		}

		cfg := new(types.VarnishConfiguration)
		cfg.Host = varnishHost
		cfg.Backends = containers

		fmt.Printf("Return varnish Config %+v \n", cfg)

		varnishConfig <- *cfg
	}

	longPoll(update)
}

func (r *Rancher) ExecuteReloadCommand(containerId string) types.ExecuteCommand{
	return executeCommand(containerId)
}