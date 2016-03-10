package service

import (
	"github.com/samalba/dockerclient"
)

type ServiceInfo struct {
	Service    Service
	Status     int
	Containers []dockerclient.ContainerInfo
}

var ServicesInfo []ServiceInfo

type Service struct {
	ContainersIds []string
	ServiceConfig ServiceConfig
}

type ServiceConfig struct {
	Kind     string          `kind`
	Metadata ServiceMetadata `metadata`
	Spec     ServiceSpec     `spec`
}

type ServiceMetadata struct {
	Name   string            `name`
	Labels map[string]string `labels`
}

type ServiceSpec struct {
	Ports     []int     `ports`
	Replicas  int       `replicas`
	Image     string    `image`
	Resources Resources `resources` //****
	Ips       []string  `ips`
	Cmd       []string  `cmd`
	Selector  []string  `selector`
}

type Resources struct {
	Memory      int64  `memory`
	MemorySwap  int64  `memory-swap`
	CpuShares   int64  `cpu-shares`
	CpusetCpus  string `cpuset-cpus`
	CpuQuota    string `cpu-quota`
	NetworkMode string `networkmode`
}
