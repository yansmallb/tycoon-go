package cli

import (
	"github.com/samalba/dockerclient"
	"github.com/yansmallb/tycoon-go/api"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
	"github.com/yansmallb/tycoon-go/swarmclient"
	"time"
)

type serviceInfo struct {
	Service    service.Service
	Status     int
	Containers []dockerclient.ContainerInfo
}

var ServicesInfo []serviceInfo

func manage(etcdPath string) error {
	// start etcd watcher
	err := etcdWatcher(ServicesInfo, etcdPath)

	// start API listener
	hosts := []string{api.TycoonHost + api.TycoonPort}
	server := api.NewServer(hosts)
	server.SetHandler(api.NewPrimary())
	err := server.ListenAndServe()
	return err
}

func etcdWatcher(servicesInfo []serviceInfo, etcdPath string) error {
	heartbeat := 20 * time.Second
	etcd, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		return err
	}
	swarm, err := swarmclient.NewSwarmClient()
	go func() {
		for {
			// get services name from etcd 
			servicesName := etcd.GetServices()
			for index := range servicesName {
				//save 
				s := etcd.GetService(servicesName[index])
				ServicesInfo[len(ServicesInfo)-1].Service = s

				containers :=  []dockerclient.ContainerInfo
				for container_index := range s.ContainersIds {
					ci,err := swarm.GetContainerInfo(s.ContainersIds[container_index])
					if(err){

					}
					containers[len(containers)-1] = ci
				}
				ServicesInfo[len(ServicesInfo)-1].Containers = s
			}
			time.Sleep(heartbeat)
			break
		}
	}()
	return err
}
