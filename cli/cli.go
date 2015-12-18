package cli

import (
	"fmt"
	"os"
)

func Run() {
	if len(os.Args) == 1 {
		help()
		return
	}
	var err error
	command := os.Args[1]
	if command == "create" {
		if len(os.Args) != 4 {
			fmt.Println("the `create` command takes two arguments. See help")
			return
		}

		filePath := os.Args[2]
		etcdPath := os.Args[3]
		err = create(filePath, etcdPath)
	}
	if command == "delete" {
		if len(os.Args) != 4 {
			fmt.Println("the `delete` command takes two arguments. See help")
			return
		}
		serviceName := os.Args[2]
		etcdPath := os.Args[3]
		err = delete(serviceName, etcdPath)
	}
	if command == "manage" {
		if len(os.Args) != 3 {
			fmt.Println("the `manage` command takes one arguments. See help")
			return
		}
		etcdPath := os.Args[2]
		err = manage(etcdPath)
	}
	if command == "help" {
		help()
	}
	if err != nil {
		fmt.Print("Error:")
		fmt.Println(err)
		return
	}
}
