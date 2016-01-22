package etcdclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/yansmallb/tycoon-go/service"
	"golang.org/x/net/context"

	"path"
	"strings"
	"time"
)

var EtcdPath = "http://127.0.0.1:2379"

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
	log.Debugf("etcdclient.GetServices():GetServices On etcd Respone %+v", Response)
	if err != nil {
		return nil, err
	}
	nodes := Response.Node.Nodes
	// Unmarshal
	servicesName := make([]string, 0)
	for index := range nodes {
		nodekey := nodes[index].Key
		nodekey = strings.Replace(nodekey, TycoonDir, "", 1)
		servicesName = append(servicesName, nodekey)
	}
	return servicesName, nil
}

func (e *Etcd) GetService(serviceName string) (*service.Service, error) {
	s := new(service.Service)
	servicePath := path.Join(TycoonDir, serviceName)

	goption := new(client.GetOptions)
	goption.Recursive = true

	// etcd get
	Response, err := e.client.Get(context.Background(), servicePath, goption)
	log.Debugf("etcdclient.GetService():GetService On etcd Respone %+v", Response)
	if err != nil {
		return s, err
	}
	nodes := Response.Node.Nodes

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
			idnodes := nodes[index].Nodes
			ids := make([]string, 0)
			for ids_index := range idnodes {
				ids = append(ids, idnodes[ids_index].Value)
			}
			s.ContainersIds = ids
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
	log.Debugf("etcdclient.DeleteService():DeleteService On etcd Respone %+v", Response)
	return err
}

func (e *Etcd) CreateService(serviceName string, serviceCfgStr string, containerIds []string) error {
	servicePath := path.Join(TycoonDir, serviceName)
	Response, err := e.client.Create(context.Background(), servicePath+"/ServiceConfig", serviceCfgStr)
	log.Debugf("etcdclient.CreateService():CreateService ServiceConfig On etcd Respone %+v", Response)
	if err != nil {
		return err
	}
	for index := range containerIds {
		Response, err := e.client.Create(context.Background(), servicePath+"/ContainerIds/"+containerIds[index], containerIds[index])
		log.Debugf("etcdclient.CreateService():CreateService ContainerIds On etcd Respone %+v", Response)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Etcd) UpdateServiceContainerIds(serviceName string, containerIds []string) error {
	servicePath := path.Join(TycoonDir, serviceName)
	doption := new(client.DeleteOptions)
	doption.Dir = true
	doption.Recursive = true

	Response, err := e.client.Delete(context.Background(), servicePath+"/ContainerIds/", doption)
	log.Debugf("etcdclient.CreateService():CreateService ServiceConfig On etcd Respone %+v", Response)
	if err != nil {
		return err
	}

	for index := range containerIds {
		Response, err := e.client.Create(context.Background(), servicePath+"/ContainerIds/"+containerIds[index], containerIds[index])
		log.Debugf("etcdclient.CreateService():CreateService ContainerIds On etcd Respone %+v", Response)
		if err != nil {
			return err
		}
	}
	return nil
}
