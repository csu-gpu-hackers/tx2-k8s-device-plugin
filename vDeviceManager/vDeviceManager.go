package vDeviceManager

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"log"
	"path"
)

type VDevice struct {
	vDeviceType string
	connection *utils.Messenger
	vlmMgr *VolumeManager
}

func (vd *VDevice) CheckStatus() utils.DeviceStatus  {
	return utils.PENDING
}

func (vd *VDevice) Serve() {
	vd.connection.Actions["Register"] = vd.Register
	go vd.connection.Serve()
}

//func (vd *VDevice) monitor() {}

func (vd *VDevice) Register(PodUID string, containerID string, registerInfo string) string {
	vd.vlmMgr.UpdateContainerID(containerID)
	vd.vlmMgr.WriteConfig()
	return "Register Success"
}

func (vd *VDevice) Report(reportSource string, reportInfo string) string {
	return "nil"
}




type VDeviceManager struct {
	vDevices []*VDevice
}

func (vdm *VDeviceManager) NewDevice(deviceType string, vlmMgr *VolumeManager)  {
	Actions := make(map[string]func(string, string, string)string)
	vd := &VDevice{
		vDeviceType: deviceType,
		vlmMgr: vlmMgr,
		connection: utils.InitMessenger(path.Join(vlmMgr.VCudaLibHostPath, "vdm.sock"), Actions),
	}
	vdm.vDevices = append(vdm.vDevices, vd)
	go vd.Serve()
	
}


func (vdm *VDeviceManager) Serve() {
	for _, vdevice := range vdm.vDevices {
		switch vdevice.CheckStatus() {
		case utils.OK:
			continue
		case utils.PENDING:
			log.Fatalf("Container Still pending after serving, please check\n")
		case utils.DEAD:
			vdevice.vlmMgr.ReleaseConfig()
		default:
			log.Println("Unexpected status")

		}
	}
}



