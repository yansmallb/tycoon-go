package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/tycoon-go/api"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
	"github.com/yansmallb/tycoon-go/swarmclient"
	"time"
)

var HeartbeatSecond = 30

func manage(etcdPath string) error {
	log.Infoln("cli.manage():Start Manage")
	// start etcd watcher
	service.ServicesInfo = make([]service.ServiceInfo, 0)
	err := servicesWatcher(etcdPath)
	if err != nil {
		log.Fatalf("cli.manage():%+v\n", err)
	}
	// start API listener
	hosts := []string{api.TycoonHost + api.TycoonPort}
	fmt.Printf("[info]:manage: hosts %+v\n", hosts)

	server := api.NewServer(hosts)
	server.SetHandler(api.NewPrimary())
	err = server.ListenAndServe()

	return err
}

func servicesWatcher(etcdPath string) error {
	log.Infoln("cli.servicesWatcher():Start servicesWatcher")
	heartbeat := time.Duration(HeartbeatSecond) * time.Second
	log.Debugf("cli.servicesWatcher(): heartbeat: %d second\n", HeartbeatSecond)
	etcd, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		log.Fatalf("cli.servicesWatcher():%+v\n", err)
		return err
	}
	swarm, err := swarmclient.NewSwarmClient()
	if err != nil {
		log.Fatalf("cli.servicesWatcher():%+v\n", err)
		return err
	}
	go func() {
		for {
			// get services name from etcd
			servicesName, err := etcd.GetServices()
			log.Infof("cli.servicesWatcher(): Services: %+v\n", servicesName)
			fmt.Printf("[info]:manage: Services %+v\n", servicesName)
			if err != nil {
				time.Sleep(heartbeat)
				continue
			}
			for index := range servicesName {
				//get service info with serviceName
				s, err := etcd.GetService(servicesName[index])
				if err != nil {
					log.Fatalf("cli.servicesWatcher():%+v\n", err)
				}

				//get containers info with containerIds
				containers, status := swarm.GetContainersInfo(s.ContainersIds)
				si := &service.ServiceInfo{Service: *s, Containers: containers, Status: status}

				service.ServicesInfo = append(service.ServicesInfo, *si)
				log.Debugf("cli.servicesWatcher(): %s info: %+v\n", servicesName[index], si)
			}
			time.Sleep(heartbeat)
		}
	}()
	return nil
}
