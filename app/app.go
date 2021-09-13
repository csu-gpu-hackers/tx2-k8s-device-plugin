package main

import (
	pluginmanager "github.com/csu-gpu-hackers/tx2-k8s-device-plugin/plugin-manager"
	log "github.com/sirupsen/logrus"
)

const (
	gpu_socket = "/var/lib/kubelet/device-plugins/csu-gpu.sock"
	gpu_path = "/home/gpu-device"
)

func main() {
	//pluginManager := pluginmanager.NewDevPlg()
	pluginManager := pluginmanager.NewPluginManager()
	go pluginManager.Run()
	log.Info("PluginManager start waiting for call")
	for {

	}

}