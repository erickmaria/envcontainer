package docker

import (
	"context"
	"errors"

	internalType "github.com/ErickMaria/envcontainer/internal/pkg/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func (docker *Docker) createNetwork(ctx context.Context, options []internalType.Network, labels map[string]string) ([]string, error) {

	getNetwork, err := docker.getNetwork(ctx, labels)
	if err != nil {
		return []string{}, err
	}

	// if len(getNetwork) != 0 {
	// 	networkIds := []string{}

	// 	// for _, net := range getNetwork {
	// 	// 	networkIds = append(networkIds, net.ID)
	// 	// }
	// }

	networkIds := []string{}
	networkIdSkip := false
	for _, netOpts := range options {

		if len(getNetwork) != 0 {
			for _, net := range getNetwork {

				if netOpts.Name == net.Name {
					networkIds = append(networkIds, net.ID)
					networkIdSkip = true
				}
			}
		}

		if networkIdSkip {
			networkIdSkip = false
			continue
		}

		if netOpts.External {

			netList, err := docker.client.NetworkList(ctx, types.NetworkListOptions{
				Filters: filters.NewArgs(
					filters.KeyValuePair{
						Key:   "name",
						Value: netOpts.Name,
					},
				),
			})

			if err != nil {
				return []string{}, err
			}

			if len(netList) == 0 {
				return []string{}, errors.New("network with name " + netOpts.Name + " does not exist")
			}

			networkIds = append(networkIds, netList[0].ID)
			continue
		}

		networkIPAMConfig := network.IPAM{}
		if netOpts.IPAM != nil {
			networkIPAMConfig.Driver = netOpts.Driver
			for _, ipamConfigOpts := range netOpts.IPAM.Config {
				networkIPAMConfig.Config = append(networkIPAMConfig.Config, network.IPAMConfig{
					Subnet:  ipamConfigOpts.Subnet,
					Gateway: ipamConfigOpts.Gateway,
				})
			}
		}

		networkConfig := types.NetworkCreate{
			CheckDuplicate: true,
			Driver:         netOpts.Driver,
			IPAM:           &networkIPAMConfig,
			Labels:         labels,
		}

		resp, err := docker.client.NetworkCreate(ctx, netOpts.Name, networkConfig)
		if err != nil {
			return []string{}, err
		}

		networkIds = append(networkIds, resp.ID)

		// fmt.Printf("Created container %s on network %s: %s\n", resp.ID, netOpts.Name, resp.Warning)
	}

	return networkIds, nil
}
