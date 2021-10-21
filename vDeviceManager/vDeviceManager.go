package vDeviceManager

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"os"
	"path"
)
const configPathPrefix string = "/home/ryan/"
type VDevice struct {
	vDeviceType string
	configPath string
	containerID string
	podUID string
	connection *utils.Messenger
}

func (vd *VDevice) CheckStatus() utils.DeviceStatus  {
	return utils.PENDING
}

func (vd *VDevice) Serve() {
	vd.connection.Actions["Register"] = vd.Register
	vd.connection.Serve()
}

func (vd *VDevice) getPodUidFromContainerID(containerID string) string {

}

func (vd *VDevice) Register(registerSource string, registerInfo string) string {
	vd.containerID = registerSource
	return "Register Success"
}

func (vd *VDevice) Report(reportSource string, reportInfo string) string {
	return "nil"
}

type VDeviceManager struct {
	vDevices []*VDevice
}

func (vdm *VDeviceManager) NewDevice(vDeviceName string, deviceType string, podUID string)  {
	configPath := path.Join(configPathPrefix, podUID)
	Actions := make(map[string]func(string, string)string)
	vd := &VDevice{
		vDeviceType: deviceType,
		configPath: configPath,
		containerID: "nil",
		podUID: podUID,
		connection: utils.InitMessenger(path.Join(configPath, "vdm.sock"), Actions),
	}
	vdm.vDevices = append(vdm.vDevices, vd)
	go vd.Serve()
	os.MkdirAll(configPath, 755)

}


func (vdm *VDeviceManager) Serve() {

}



