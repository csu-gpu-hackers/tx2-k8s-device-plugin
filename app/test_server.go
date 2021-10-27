package main

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	log "github.com/sirupsen/logrus"
)

func respond(source string, info string) string {
	log.Println("Received from " + source + info)
	return "GOOD!"
}


func main() {
	actions := make(map[string]func(string, string)string)
	actions["1"] = respond
	messenger := utils.InitMessenger("test.sock", actions)
	messenger.Serve()
}

//{"MessageType":"1","ContainerID":"","PodUID":"wdfghhjhgfdfghjhgf","MessageContent":"aaa"}
