package cli

import (
	"fmt"
)

func help() {
	var helpString = `
  Usage: tycoon COMMAND [args...]
  Version: 0.1.0
  Author:yansmallb
  Email:yanxiaoben@iie.ac.cn

  Commands:
      create
      delete
      manage
      help
  `
	fmt.Println(helpString)
}
