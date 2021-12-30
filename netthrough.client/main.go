package main

import (
	goutils "github.com/typa01/go-utils"
	"netthrough.client/business"
)

func main() {
	//生成客户端ID
	guid := goutils.GUID()
	task := &business.Task{
		ClientId: guid,
	}
	defer task.Stop()
	task.Start()
	select {}

}
