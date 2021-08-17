package main

import (
	device_plugin "dev-play/device-plugin"
	"dev-play/utils"
	log "github.com/sirupsen/logrus"
)

const (
	gpu_socket = "/var/lib/kubelet/device-plugins/csu-gpu.sock"
	gpu_path = "/home/ryan/gpu-device"
)

func main() {
	devPlg := device_plugin.NewDevPlg("csu.ac.cn/gpu", gpu_socket)
	log.Println("construction of dp finished, start running")
	go devPlg.Run()
	go func() {
		err := devPlg.RegisterToKubelet()
		utils.Check(err)
		log.Println("Register finished")
	}()
	for {
		
	}





}