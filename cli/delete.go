package cli

import (
	//"fmt"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
)

func delete(serviceName string, etcdPath string) error {
	client, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		return err
	}

	//quary service
	s, err := client.GetService(serviceName)
	if err != nil {
		return err
	}

	//delete from etcd
	err = client.DeleteService(serviceName)
	if err != nil {
		return err
	}

	//delete from swarm
	err = service.DeleteService(s)
	if err != nil {
		return err
	}
	return nil
}
