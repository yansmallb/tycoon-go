package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/yansmallb/tycoon-go/etcdclient"
	"os"
)

func Run() {
	if len(os.Args) == 1 {
		help()
		return
	}
	var err error
	command := os.Args[1]
	log.Debugf("cli.Run(): cli args:%+v\n", os.Args)
	if command == "create" {
		if len(os.Args) != 4 {
			createErr := "the `create` command takes two arguments. See help"
			fmt.Println(createErr)
			log.Errorln(createErr)
			return
		}
		filePath := os.Args[2]
		etcdPath := os.Args[3]
		err = create(filePath, etcdPath)
	}
	if command == "delete" {
		if len(os.Args) != 4 {
			deleteErr := "the `delete` command takes two arguments. See help"
			fmt.Println(deleteErr)
			log.Errorln(deleteErr)
			return
		}
		serviceName := os.Args[2]
		etcdPath := os.Args[3]
		err = delete(serviceName, etcdPath)
	}
	if command == "manage" {
		if len(os.Args) > 3 {
			manageErr := "the `manage` command takes one argument at most. See help"
			fmt.Println(manageErr)
			log.Errorln(manageErr)
			return
		}
		//etcdPath := os.Args[2]
		//etcdclient.EtcdPath = etcdPath
		err = manage(etcdclient.EtcdPath)
	}
	if command == "help" {
		help()
	}
	if err != nil {
		fmt.Print("[Error]:")
		fmt.Println(err)
		log.Fatalf("cli.Run():%+v\n", err)
		return
	}
}
