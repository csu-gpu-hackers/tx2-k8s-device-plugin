package devices

import (
	"github.com/fsnotify/fsnotify"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"tx2-k8s-device-plugin/utils"
)

type DeviceHandler struct{
	deviceType string
}

func (dh *DeviceHandler) GetDeviceParts() []*plugin.Device {
	panic("implement me")
}

func (dh *DeviceHandler) GetDeviceStatus() utils.DeviceStatus {
	panic("implement me")
}

func (dh *DeviceHandler) GetDeviceLoads() int {
	panic("implement me")
}

func (dh *DeviceHandler) WatchDevice() error {
	panic("implement me")
}

func (dh *DeviceHandler) GetDeviceChangeNotifier() *fsnotify.Watcher {
	panic("implement me")
}

