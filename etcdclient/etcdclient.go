package etcdclient

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/yansmallb/tycoon-go/service"
	"golang.org/x/net/context"
	//"net/http"
	"path"
	//"strconv"
	//"encoding/json"
	"strings"
	"time"
)

type Etcd struct {
	client client.KeysAPI
}

// ten years
var ServcieTimeout = 10 * 24 * 365 * time.Hour
var TycoonDir = "/yansmallb/tycoon/services/"

func NewEtcdClient(etcdpath string) (*Etcd, error) {
	endpoints := strings.Split(etcdpath, ",")
	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavctxailable
		HeaderTimeoutPerRequest: ServcieTimeout,
	}
	c, err := client.New(cfg)
	kapi := client.NewKeysAPI(c)

	etcd := new(Etcd)
	etcd.client = kapi
	return etcd, err
}

func (e *Etcd) GetServices() ([]string, error) {
	servicePath := path.Join(TycoonDir)
	goption := new(client.GetOptions)
	goption.Recursive = true

	Response, err := e.client.Get(context.Background(), servicePath, goption)
	if err != nil {
		return s, err
	}
	nodes := Response.Node.Nodes
	// Unmarshal
	servicesName := new([]string)
	for index := range nodes {
		servicesName[len(serviceName)-1] = nodes[index].Key
	}
	return servicesName
}

func (e *Etcd) GetService(serviceName string) (*service.Service, error) {
	s := new(service.Service)
	servicePath := path.Join(TycoonDir, serviceName)
	goption := new(client.GetOptions)
	goption.Recursive = true

	// etcd get
	Response, err := e.client.Get(context.Background(), servicePath, goption)
	if err != nil {
		return s, err
	}
	nodes := Response.Node.Nodes

	//fmt.Println("[Info]EtcdClient.GetService:")

	// Unmarshal
	for index := range nodes {
		if nodes[index].Key == servicePath+"/ServiceConfig" {
			sc, err := service.UnmarshalYaml([]byte(nodes[index].Value))
			if err != nil {
				return s, err
			}
			s.ServiceConfig = *sc
		}
		if nodes[index].Key == servicePath+"/ContainerIds" {
			ipnodes := nodes[index].Nodes
			ips := make([]string, 100)
			for ips_index := range ipnodes {
				ips[len(ips)-1] = ipnodes[ips_index].Value
			}
			s.ContainersIds = ips
		}
	}

	return s, err
}

func (e *Etcd) DeleteService(serviceName string) error {
	servicePath := path.Join(TycoonDir, serviceName)
	doption := new(client.DeleteOptions)
	doption.Dir = true
	doption.Recursive = true

	Response, err := e.client.Delete(context.Background(), servicePath, doption)
	fmt.Println(Response)
	return err
}

func (e *Etcd) CreateService(serviceName string, serviceCfgStr string, containerIds []string) error {
	servicePath := path.Join(TycoonDir, serviceName)
	Response, err := e.client.Create(context.Background(), servicePath+"/ServiceConfig", serviceCfgStr)
	fmt.Println(Response)
	if err != nil {
		return err
	}
	for index := range containerIds {
		Response, err := e.client.Create(context.Background(), servicePath+"/ContainerIds/"+containerIds[index], containerIds[index])
		fmt.Println(Response)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
func CreateService(c client.Client, config *service.ServiceConfig) error {
	servicePath := TycoonDir + config.Metadata.Name
	kapi := client.NewKeysAPI(c)
	//ServiceConfig
	kapi.Create(context.Background(), servicePath+"/Kind", config.Kind)
	//Metadata
	kapi.Create(context.Background(), servicePath+"/Metadata/Name", config.Metadata.Name)
	for label := range config.Metadata.Labels {
		//kapi.Create(context.Background(), servicePath+"/Metadata/Label/"+label, label)
		fmt.Println(label)
	}
	//Spec
	for port := range config.Spec.Ports {
		kapi.Create(context.Background(), servicePath+"/Spec/Ports/"+strconv.Itoa(config.Spec.Ports[port]), strconv.Itoa(config.Spec.Ports[port]))
	}
	kapi.Create(context.Background(), servicePath+"/Spec/Replicas", strconv.Itoa(config.Spec.Replicas))
	kapi.Create(context.Background(), servicePath+"/Spec/Image", config.Spec.Image)
	for resource := range config.Spec.Resources {
		//kapi.Create(context.Background(), servicePath+"/Spec/Resources/"+resource, resource)
		fmt.Println(resource)
	}
	for ip := range config.Spec.Ips {
		kapi.Create(context.Background(), servicePath+"/Spec/Ips/"+config.Spec.Ips[ip], config.Spec.Ips[ip])
	}
	for selector := range config.Spec.Selector {
		//kapi.Create(context.Background(), servicePath+"/Spec/Selector/"+selector, selector)
		fmt.Println(selector)
	}
	return nil
}
*/

func (e *Etcd) WatchOnce() ([]string, error) {
	s := "1,2,3"
	res := strings.Split(s, ",")
	return res, nil
}
