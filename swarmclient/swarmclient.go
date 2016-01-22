package swarmclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

var SwarmHost = "http://127.0.0.1"
var SwarmPort = ":2375"

type SwarmClient struct {
	Client *dockerclient.DockerClient
}

func NewSwarmClient() (*SwarmClient, error) {
	swarmUrl := SwarmHost + SwarmPort
	docker, err := dockerclient.NewDockerClient(swarmUrl, nil)
	sc := &SwarmClient{Client: docker}
	return sc, err
}

func (sc *SwarmClient) CreateContainer(config *dockerclient.ContainerConfig, name string) (string, error) {
	return sc.Client.CreateContainer(config, name)
}

func (sc *SwarmClient) StartContainer(id string, config *dockerclient.HostConfig) error {
	return sc.Client.StartContainer(id, config)
}

func (sc *SwarmClient) StopContainer(id string, timeout int) error {
	return sc.Client.StopContainer(id, timeout)
}

func (sc *SwarmClient) RestartContainer(id string, timeout int) error {
	return sc.Client.RestartContainer(id, timeout)
}

func (sc *SwarmClient) RemoveContainer(id string, force, volumes bool) error {
	return sc.Client.RemoveContainer(id, force, volumes)
}

func (sc *SwarmClient) GetContainerInfo(containerId string) (*dockerclient.ContainerInfo, error) {
	return sc.Client.InspectContainer(containerId)
}

func (sc *SwarmClient) GetContainersInfo(containersIds []string) ([]dockerclient.ContainerInfo, int) {
	status := 1
	containers := make([]dockerclient.ContainerInfo, 0)
	for _, container := range containersIds {
		ci, err := sc.GetContainerInfo(container)
		if err != nil {
			status = -1
			continue
		}
		if ci.State.Running != true {
			status = -1
		}
		containers = append(containers, *ci)
	}
	return containers, status
}

func (sc *SwarmClient) ContainersHealthCheck(containers []dockerclient.ContainerInfo) int {
	status := 1
	for _, container := range containers {
		if container.State.Running {
			continue
		}
		status = -1
		// try to restart
		err := sc.RestartContainer(container.Id, 0)
		if err != nil {
			log.Fatalf("swarmclient.ContainersHealthCheck():%s\n", err)
			err = sc.RemoveContainer(container.Id, true, false)
			if err != nil {
				log.Debugf("swarmclient.ContainersHealthCheck():%s\n", err)
			}
		} else {
			status = 1
		}
	}
	return status
}
