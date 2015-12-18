package swarmclient

import (
	"github.com/samalba/dockerclient"
	//"strconv"
)

var swarmhost = "http://127.0.0.1"
var swarmport = ":2375"

func NewSwarmClient() (*dockerclient.DockerClient, error) {
	swarmUrl := swarmhost + swarmport
	docker, err := dockerclient.NewDockerClient(swarmUrl, nil)
	return docker, err
}

func (dc *dockerclient.DockerClient) GetContainerInfo(containerId string) (*ContainerInfo, error) {
	return dc.InspectContainer(containerId)
}
