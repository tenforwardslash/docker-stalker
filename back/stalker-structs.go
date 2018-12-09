package main

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types"
)

func GetStalkerPorts(ports []types.Port) []*StalkerPort {
	var containerPorts []*StalkerPort
	//loop and create
	for _, p := range ports {
		containerPorts = append(containerPorts, &StalkerPort{
			Private: p.PrivatePort,
			Public: p.PublicPort,
			Type: p.Type,
		})
	}

	return containerPorts
}

func GetStalkerMounts(mounts []types.MountPoint) []*StalkerMount {
	var containerMounts []*StalkerMount
	//loop and create
	for _, m := range mounts {
		containerMounts = append(containerMounts, &StalkerMount{
			Type: m.Type,
			Source: m.Source,
			Destination: m.Destination,
		})
	}

	return containerMounts
}

type StalkerContainer struct {
	Name        string          `json:"name"`
	Image       string          `json:"image"`
	Created     int64           `json:"created"`
	Status      string          `json:"status"`
	State       string          `json:"state"`
	Ports       []*StalkerPort  `json:"ports"`
	Mounts      []*StalkerMount `json:"mounts"`
	EnvVars     []string        `json:"envVars"`
	Network     string          `json:"network"`
	ContainerId string          `json:"containerId"`
}

type StalkerPort struct {
	Private uint16 `json:"private"`
	Public  uint16 `json:"public"`
	Type    string `json:"type"`
}

type StalkerMount struct {
	Type        mount.Type `json:"type"`
	Source      string     `json:"source"`
	Destination string     `json:"destination"`
}