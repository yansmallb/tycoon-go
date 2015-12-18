package cli

import (
	"fmt"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
	"io/ioutil"
	"path/filepath"
)

func create(localfilepath string, etcdpath string) error {
	filename, _ := filepath.Abs(localfilepath)
	yamlconfig, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// unmarshal yaml file
	config, err := service.UnmarshalYaml(yamlconfig)
	if err != nil {
		fmt.Println(err)
		//return err
	}
	fmt.Printf("Value: %#v\n", config)
	//create service on swarm
	//containerIds, err := service.CreateService(config)
	containerIds := []string{"A", "B"}

	//create service on etcd
	client, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		return err
	}
	err = client.CreateService(config.Metadata.Name, string(yamlconfig), containerIds)
	if err != nil {
		return err
	}
	return nil
}
