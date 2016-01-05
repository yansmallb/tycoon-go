package cli

import (
	"fmt"
)

func help() {
	var helpString = `
  Usage: tycoon COMMAND [args...]
  Version: 0.2.0
  Author:yansmallb
  Email:yanxiaoben@iie.ac.cn

  Commands:
      create    [localyaml path] [etcd path]
      delete    [service name]  [etcd path]
      manage  [etcdpath]
      help
  `
	fmt.Println(helpString)
}
