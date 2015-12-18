package service

import (
	"fmt"
	"github.com/samalba/dockerclient"
	"github.com/yansmallb/tycoon-go/swarmclient"
	"gopkg.in/yaml.v2"
	"strconv"
)

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
	Name   string
	Labels map[string]string
}
type ServiceSpec struct {
	Ports     []int
	Replicas  int
	Image     string
	Resources []string //****
	Ips       []string
	Selector  []string
}

func UnmarshalYaml(in []byte) (*ServiceConfig, error) {
	config := new(ServiceConfig)
	err := yaml.Unmarshal(in, &config)
	return config, err
}

func CreateService(config *ServiceConfig) ([]string, error) {
	swarm, err := swarmclient.NewSwarmClient()
	containerIds := make([]string, 100)
	containerConfig := new(dockerclient.ContainerConfig)
	containerConfig.Image = config.Spec.Image
	containerConfig.Labels = config.Metadata.Labels
	portBindings := make(map[string][]dockerclient.PortBinding)
	containerConfig.Env = config.Spec.Selector

	if config.Spec.Replicas != 0 {
		for i := 0; i < config.Spec.Replicas; i++ {
			for index := range config.Spec.Ports {
				portbinding := &dockerclient.PortBinding{HostIp: "0.0.0.0", HostPort: strconv.Itoa(config.Spec.Ports[index])}
				portBindings[strconv.Itoa(index)][0] = *portbinding
			}
			hostConfig := &dockerclient.HostConfig{PortBindings: portBindings}
			containerConfig.HostConfig = *hostConfig
			containerId, err := swarm.CreateContainer(containerConfig, "")
			containerIds[len(containerIds)-1] = containerId
			fmt.Println(err)
		}
	} else {
		for ip_index := range config.Spec.Ips {
			for port_index := range config.Spec.Ports {
				portbinding := &dockerclient.PortBinding{HostIp: config.Spec.Ips[ip_index], HostPort: strconv.Itoa(config.Spec.Ports[port_index])}
				portBindings[strconv.Itoa(port_index)][0] = *portbinding
			}
			hostConfig := &dockerclient.HostConfig{PortBindings: portBindings}
			containerConfig.HostConfig = *hostConfig
			containerId, err := swarm.CreateContainer(containerConfig, "")
			containerIds[len(containerIds)-1] = containerId
			fmt.Println(err)
		}
	}
	return containerIds, err
}
func DeleteService(s *Service) error {
	swarm, err := swarmclient.NewSwarmClient()
	for index := range s.ContainersIds {
		swarm.StopContainer(s.ContainersIds[index], 20)
	}
	return err
}

/*
func GetServices() (services Service,error){

	return nil,nil
}
*/
