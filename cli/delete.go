package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
)

func delete(serviceName string, etcdPath string) error {
	log.Infoln("cli.delete():Start Delete")
	client, err := etcdclient.NewEtcdClient(etcdPath)
	if err != nil {
		log.Fatalf("cli.delete():%+v\n", err)
		fmt.Printf("[error]cli.delete():%+v\n", err)
		return err
	}

	//quary service
	s, err := client.GetService(serviceName)
	if err != nil {
		log.Fatalf("cli.delete():%+v\n", err)
		fmt.Printf("[error]cli.delete():%+v\n", err)
		return err
	}

	//delete from etcd
	err = client.DeleteService(serviceName)
	if err != nil {
		log.Fatalf("cli.delete():%+v\n", err)
		fmt.Printf("[error]cli.delete():%+v\n", err)
	}

	//delete from swarm
	err = service.DeleteService(s)
	if err != nil {
		log.Fatalf("cli.delete():%+v\n", err)
		fmt.Printf("[error]cli.delete():%+v\n", err)
		return err
	}
	return nil
}
