package devices

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"github.com/fsnotify/fsnotify"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type DeviceHandler struct{
	deviceName   string
	DevicePath   string
	deviceSocket string
	DeviceParts []*plugin.Device
	DeviceChangeNotifier *fsnotify.Watcher
}

func NewDeviceHandler(devicePartsNum int, deviceName string, deviceSocket string) *DeviceHandler {
	deviceParts := make([]*plugin.Device, devicePartsNum)
	for i := 0; i < devicePartsNum; i++ {
		deviceParts[i] = &plugin.Device{
			ID:                   string(i),
			Health:               plugin.Healthy,
			Topology:             nil,
		}
	}
	DeviceChangeNotifier, err := fsnotify.NewWatcher()
	utils.Check(err)
	return &DeviceHandler{
		deviceName:           deviceName,
		deviceSocket:         deviceSocket,
		DeviceParts:          deviceParts,
		DeviceChangeNotifier: DeviceChangeNotifier,
	}
}

func (dh *DeviceHandler) GetDeviceParts() []*plugin.Device {
	return dh.DeviceParts
}

func (dh *DeviceHandler) GetDeviceStatus() utils.DeviceStatus {
	//if h.GetDeviceLoads() > 95 {
	//	return utils.OCCUPIED
	//} else {
	//	return utils.OK
	//}
	return utils.OK
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

