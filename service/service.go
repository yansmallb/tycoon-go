package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/yansmallb/tycoon-go/swarmclient"
	"gopkg.in/yaml.v2"
	"strconv"
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
	Ports     []int    `ports`
	Replicas  int      `replicas`
	Image     string   `image`
	Resources []string //****
	Ips       []string `ips`
	Cmd       []string `cmd`
	Selector  []string `selector`
}

func UnmarshalYaml(in []byte) (*ServiceConfig, error) {
	config := new(ServiceConfig)
	err := yaml.Unmarshal(in, &config)
	return config, err
}

func MarshalYaml(sc *ServiceConfig) (string, error) {
	out, err := yaml.Marshal(sc)
	return string(out), err
}

type ServiceFunc func(*Service) error

func CreateService(config *ServiceConfig) ([]string, error) {
	swarm, err := swarmclient.NewSwarmClient()

	if err != nil {
		log.Fatalf("service.CreateService():%s\n", err)
		return nil, err
	}
	containerIds := make([]string, 0)
	containerConfig := new(dockerclient.ContainerConfig)
	containerConfig.Image = config.Spec.Image
	containerConfig.Labels = config.Metadata.Labels
	portBindings := make(map[string][]dockerclient.PortBinding)
	containerConfig.Env = config.Spec.Selector
	containerConfig.Cmd = config.Spec.Cmd

	//exposed ports , so that others can't use
	if len(config.Spec.Ports) > 0 {
		ports := make(map[string]struct{})
		for index := range config.Spec.Ports {
			port := strconv.Itoa(config.Spec.Ports[index])
			ports[port] = struct{}{}
		}
		containerConfig.ExposedPorts = ports
	}

	// use host, create and start containers
	hostConfig := &dockerclient.HostConfig{NetworkMode: "host"}
	numOfTimes := 0

	if config.Spec.Replicas > 0 {
		numOfTimes = config.Spec.Replicas
	} else {
		numOfTimes = len(config.Spec.Ips)
		if len(config.Spec.Ports) == 0 {
			err := fmt.Errorf("service.CreateService():%+s Give ips but not give ports\n", config.Metadata.Name)
			log.Error(err)
			return nil, err
		}
	}

	for i := 0; i < numOfTimes; i++ {
		// use to filter specific ips and ports
		for _, port := range config.Spec.Ports {
			portbinding := &dockerclient.PortBinding{HostPort: strconv.Itoa(port)}
			if config.Spec.Replicas <= 0 {
				portbinding.HostIp = config.Spec.Ips[i]
			}
			portBindings[strconv.Itoa(port)] = []dockerclient.PortBinding{*portbinding}
		}
		hostConfig.PortBindings = portBindings
		containerConfig.HostConfig = *hostConfig

		// give container different name
		containerName := config.Metadata.Name
		if role := config.Metadata.Labels["role"]; role != "" {
			containerName += "_" + role
		}
		containerName += strconv.Itoa(i)
		log.Debugf("service.CreateService():containerName:%s ; containerConfig:%+v\n", containerName, containerConfig)

		//create container
		containerId, err := swarm.CreateContainer(containerConfig, containerName)
		log.Debugf("service.CreateService():containerId:%s\n", containerId)
		if err != nil {
			log.Fatalf("service.CreateService():%s\n", err)
			return nil, err
		}
		//start container
		swarm.StartContainer(containerId, hostConfig)
		if err != nil {
			log.Fatalf("service.CreateService():%s\n", err)
			return nil, err
		}
		containerIds = append(containerIds, containerId)
	}
	return containerIds, err
}

func DeleteService(s *Service) error {
	swarm, err := swarmclient.NewSwarmClient()
	if err != nil {
		return err
	}
	log.Debugf("service.DeleteService():containerIds:%v\n", s.ContainersIds)
	for index := range s.ContainersIds {
		err := swarm.StopContainer(s.ContainersIds[index], 0)
		if err != nil {
			log.Errorf("service.DeleteService()::%v\n", err)
		}
		err = swarm.RemoveContainer(s.ContainersIds[index], true, false)
		if err != nil {
			log.Errorf("service.DeleteService()::%v\n", err)
		}
	}
	return nil
}

func GetService(s *Service) (*ServiceInfo, error) {
	swarm, err := swarmclient.NewSwarmClient()
	if err != nil {
		return nil, err
	}
	containers, status := swarm.GetContainersInfo(s.ContainersIds)
	if status != 1 {
		err = fmt.Errorf("containers status is not good")
		log.Errorf("service.GetService():%s:%s", s.ServiceConfig.Metadata.Name, err)
	}
	si := &ServiceInfo{Service: *s,
		Status:     status,
		Containers: containers}
	return si, err
}

// use for web API
func RestartService(s *Service) error {
	swarm, err := swarmclient.NewSwarmClient()
	if err != nil {
		log.Errorf("service.RestartService():%v\n", err)
		return err
	}
	for index := range s.ContainersIds {
		err_for := swarm.StopContainer(s.ContainersIds[index], 0)
		if err_for != nil {
			err = err_for
			log.Errorf("service.RestartService():%v\n", err)
		}
		err_for = swarm.StartContainer(s.ContainersIds[index], nil)
		if err_for != nil {
			err = err_for
			log.Errorf("service.RestartService():%v\n", err)
		}
	}
	return err
}

//use for web API ,Best Not Use
func StopService(s *Service) error {
	swarm, err := swarmclient.NewSwarmClient()
	if err != nil {
		log.Errorf("service.StopService():%v\n", err)
		return err
	}
	for index := range s.ContainersIds {
		err_for := swarm.StopContainer(s.ContainersIds[index], 0)
		if err_for != nil {
			err = err_for
			log.Errorf("service.StopService():%v\n", err)
		}
	}
	return err
}

//use for web API ,Best Not Use
func StartService(s *Service) error {
	swarm, err := swarmclient.NewSwarmClient()
	if err != nil {
		log.Errorf("service.StartService():%v\n", err)
		return err
	}
	for index := range s.ContainersIds {
		err_for := swarm.StartContainer(s.ContainersIds[index], nil)
		if err_for != nil {
			err = err_for
			log.Errorf("service.StartService():%v\n", err)
		}
	}
	return err
}
