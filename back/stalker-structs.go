package main

import (
	"github.com/docker/docker/api/types"
)

type StalkerContainer struct {
	Name        string             `json:"name"`
	Image       string             `json:"image"`
	Created     int64              `json:"created"`
	Status      string             `json:"status"`
	State       string             `json:"state"`
	Ports       []types.Port       `json:"ports"`
	Mounts      []types.MountPoint `json:"mounts"`
	EnvVars     []string           `json:"envVars"`
	Network     string             `json:"network"`
	ContainerId string             `json:"containerId"`
}