package main

import (
	log "github.com/sirupsen/logrus"
	pluginmanager "tx2-k8s-device-plugin/plugin-manager"
)

const (
	gpu_socket = "/var/lib/kubelet/device-plugins/csu-gpu.sock"
	gpu_path = "/home/gpu-device"
)

func main() {
	pluginManager := pluginmanager.NewPluginManager()
	go pluginManager.Run()
	log.Info("PluginManager start waiting for call")
	for {

	}

}