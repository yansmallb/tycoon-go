package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
	"io/ioutil"
	"path/filepath"
)

func create(localfilepath string, etcdpath string) error {
	log.Infoln("cli.create():Start Create")
	log.Debugln("cli.create():ConfigFile Path:" + localfilepath)
	filename, _ := filepath.Abs(localfilepath)
	yamlconfig, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cli.create():%+v\n", err)
		return err
	}
	// unmarshal yaml file
	config, err := service.UnmarshalYaml(yamlconfig)
	if err != nil {
		log.Fatalf("cli.create():%+v\n", err)
		return err
	}
	fmt.Print("[Info] yaml config: ")
	fmt.Println(config)
	log.Debugf("cli.create():yaml config:%+v\n", config)

	//create service on swarm
	containerIds, err := service.CreateService(config)
	if err != nil {
		return err
	}
	fmt.Print("[Info] containerIDs: ")
	fmt.Println(containerIds)
	log.Infof("cli.create():containerIDs :%+v\n", containerIds)

	//create service on etcd
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		log.Fatalf("cli.create():%+v\n", err)
		return err
	}
	err = client.CreateService(config.Metadata.Name, string(yamlconfig), containerIds)
	if err != nil {
		log.Fatalf("cli.create():%+v\n", err)
		return err
	}
	return nil
}
