package main

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"log"
)

func main() {
	poduid := "85a221ee-f9d5-489a-b4a9-20dee177497a"
	log.Println(utils.CheckPodStatus(poduid))
}
