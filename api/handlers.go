package api

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/service"
	"io/ioutil"
	"net/http"
)

func getServices(w http.ResponseWriter, r *http.Request) {
	// servicesName save on etcd
	etcd, err := etcdclient.NewEtcdClient(etcdclient.EtcdPath)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	servicesName, err := etcd.GetServices()
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(servicesName)
}

func getServicesInfo(w http.ResponseWriter, r *http.Request) {
	// return manage servicesInfo
	// because manage watch the servicesInfo every heartbeat
	// if we getServicesInfo from swarm ,there are many cost
	json.NewEncoder(w).Encode(service.ServicesInfo)
}

func getService(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]
	etcd, err := etcdclient.NewEtcdClient(etcdclient.EtcdPath)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s, err := etcd.GetService(serviceName)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	si, err := service.GetService(s)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(si)
}

//****
func createService(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonconfig, _ := ioutil.ReadAll(r.Body)
	config := new(service.ServiceConfig)
	json.Unmarshal([]byte(jsonconfig), &config)

	//create service on swarm
	containerIds, err := service.CreateService(config)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Print("[Info] containerIDs ")
	//fmt.Println(containerIds)

	//create service on etcd
	client, err := etcdclient.NewEtcdClient(etcdclient.EtcdPath)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	yamlconfig, err := service.MarshalYaml(config)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = client.CreateService(config.Metadata.Name, string(yamlconfig), containerIds)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("CreateService" + " Service [" + config.Metadata.Name + "] success")
}

func deleteService(w http.ResponseWriter, r *http.Request) {
	doService(w, r, service.DeleteService, "DeleteService")
}

func restartService(w http.ResponseWriter, r *http.Request) {
	doService(w, r, service.RestartService, "RestartService")
}

func stopService(w http.ResponseWriter, r *http.Request) {
	doService(w, r, service.StopService, "StopService")
}

func startService(w http.ResponseWriter, r *http.Request) {
	doService(w, r, service.StartService, "StartService")
}

// for start,stop,restart,delete
func doService(w http.ResponseWriter, r *http.Request, commond service.ServiceFunc, comstr string) {
	serviceName := mux.Vars(r)["name"]
	etcd, err := etcdclient.NewEtcdClient(etcdclient.EtcdPath)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s, err := etcd.GetService(serviceName)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = commond(s)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		json.NewEncoder(w).Encode(comstr + " Service [" + serviceName + "] success")
	}
}

func httpError(w http.ResponseWriter, err string, status int) {
	log.WithField("status", status).Errorf("api.httpError():HTTP error: %v", err)
	http.Error(w, err, status)
}
