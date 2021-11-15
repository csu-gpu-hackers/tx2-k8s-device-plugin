package main

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/device-plugin"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	gpu_socket = "/var/lib/kubelet/device-plugins/csu-gpu.sock"
	gpu_path = "/home/gpu-device"
)

func main() {
	_, err := os.Stat(gpu_socket)
	if err == nil {
		err := os.Remove(gpu_socket)
		utils.Check(err)
	}

	devPlg := device_plugin.NewDevPlg("csu.ac.cn/gpu", gpu_socket)
	log.Println("construction of dp finished, start running")
	go devPlg.Run()
	go func() {
		err := devPlg.RegisterToKubelet()
		utils.Check(err)
		log.Println("Register finished")
	}()
	select {}

}
