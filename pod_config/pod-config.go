package pod_config

import (
	"path"
)

var ConfigBasePath string
type Status int8
const (
	OK Status = iota
	DEAD
	INITIALIZING
)

type ContainerConfig struct {
	ContainerID     string
	ConfigPath      string
	ContainerStatus Status
}

func NewContainer(ContainerID string) *ContainerConfig{
	return &ContainerConfig{
		ContainerID:     ContainerID,
		ConfigPath:      path.Join(ConfigBasePath, ContainerID),
		ContainerStatus: INITIALIZING,
	}
}

