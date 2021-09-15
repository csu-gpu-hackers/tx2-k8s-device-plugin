package devices

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type GPUManager struct {
	deviceName   string
	DevicePath   string
	deviceSocket string
	DeviceParts []*plugin.Device
	DeviceChangeNotifier *fsnotify.Watcher
	//DeviceChangeNotifier chan bool
}

func NewGPUManager() *GPUManager {
	DevicePercents := make([]*plugin.Device, 100)
	for i := 0; i < 100; i++ {
		DevicePercents[i] = &plugin.Device{
			ID:                   string(i),
			Health:               plugin.Healthy,
			Topology:             nil,
		}
	}

	DeviceChangeNotifier, err := fsnotify.NewWatcher()
	utils.Check(err)
	//defer DeviceChangeNotifier.Close()

	return &GPUManager{
		deviceName:   	"xwan-gpu",
		//DevicePath:   "/sys/devices/platform/host1x/17000000.gp10b/load",
		DevicePath:		"/home/gpu-device",
		deviceSocket: 	"xwan-gpu.sock",
		DeviceParts:  	DevicePercents,
		DeviceChangeNotifier: DeviceChangeNotifier,
		//DeviceChangeNotifier: make(chan bool),
	}
}

//func (gpu *GPUManager) Allocate()  {
//}

func (gpu *GPUManager) GetDeviceStatus() utils.DeviceStatus {
	if gpu.GetDeviceLoads() > 95 {
		return utils.OCCUPIED
	} else {
		return utils.OK
	}
}

func (gpu *GPUManager) GetDeviceLoads() int {
	return gpu.updateDeviceLoads()
}

func (gpu *GPUManager) GetDeviceParts() []*plugin.Device {
	return gpu.DeviceParts
}

func (gpu *GPUManager) updateDeviceLoads() int {
	load, err := utils.ExtractNumber(utils.ReadFile(gpu.DevicePath))
	utils.Check(err)
	//load = load / 10

	//for i := 0; i < 100; i++ {
	//	gpu.DeviceParts[i].Health = plugin.Healthy
	//}
	//
	//for i := 0; i < load; i++ {
	//	gpu.DeviceParts[i].Health = plugin.Unhealthy
	//}
	return load
}

func (gpu *GPUManager) GetDeviceChangeNotifier() *fsnotify.Watcher {
	return gpu.DeviceChangeNotifier
}

func (gpu *GPUManager) WatchDevice() error {
	log.Println("Watch Device start working")
	err := gpu.DeviceChangeNotifier.Add(gpu.DevicePath)
	log.Println("watching device: ", gpu.DevicePath)
	utils.Check(err)
	return nil
}
