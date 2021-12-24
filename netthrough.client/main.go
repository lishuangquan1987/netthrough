package main

import (
	"fmt"

	"netthrough.client/config"
)

func main() {
	fmt.Println(config.Config.ServerIp)
	fmt.Println(config.Config.ServerPort)
}
