package main

import (
	"fmt"
	"github.com/yansmallb/tycoon-go/api"
	"github.com/yansmallb/tycoon-go/cli"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"github.com/yansmallb/tycoon-go/swarmclient"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	TycoonHost string `TycoonHost`
	TycoonPort string `TycoonPort`
	SwarmHost  string `SwarmHost`
	SwarmPort  string `SwarmPort`
	EtcdPath   string `EtcdPath`
	LogPath    string `LogPath`
	LogLevel   string `LogLevel`
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	file, _ := exec.LookPath(os.Args[0])
	p, _ := filepath.Abs(file)
	yamlpath := path.Join(p, "../config/config.yaml")
	err := UnmarshalConfig(yamlpath)
	if err != nil {
		fmt.Print("main.main():[Error]:")
		fmt.Println(err)
		log.Fatalf("main.main():%+v\n", err)
		return
	}
	log.Infoln("main.main():Start Tycoon")
	cli.Run()
}

func UnmarshalConfig(path string) error {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	c := new(Config)
	err = yaml.Unmarshal(in, &c)
	if err != nil {
		return err
	}

	if c.LogPath != "" {
		p, _ := filepath.Abs(c.LogPath)
		if !checkFileIsExist(p) {
			os.Create(p)
		}
		file, err := os.OpenFile(p, os.O_RDWR|os.O_APPEND|os.O_CREATE, 777)
		if err != nil {
			return err
		}
		log.SetOutput(file)
		log.Infoln("main.UnmarshalConfig():Log Path : " + p)
	} else {
		log.Infoln("main.UnmarshalConfig():Log Path : os.stderr")
	}
	if c.LogLevel != "" {
		level, err := log.ParseLevel(c.LogLevel)
		if err != nil {
			log.Errorln(err)
		} else {
			log.SetLevel(level)
			log.Infoln("main.UnmarshalConfig():Log Level : " + c.LogLevel)
		}
	} else {
		log.Infoln("main.UnmarshalConfig():Log Level : debug")
	}

	api.TycoonHost = c.TycoonHost
	api.TycoonPort = c.TycoonPort
	log.Infoln("main.UnmarshalConfig():Tycoon Host : " + api.TycoonHost + api.TycoonPort)

	swarmclient.SwarmHost = c.SwarmHost
	swarmclient.SwarmPort = c.SwarmPort
	log.Infoln("main.UnmarshalConfig():Swarm Host : " + swarmclient.SwarmHost + swarmclient.SwarmPort)

	etcdclient.EtcdPath = c.EtcdPath
	log.Infoln("main.UnmarshalConfig():Etcd Host : " + etcdclient.EtcdPath)
	return nil
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
