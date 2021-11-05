package vDeviceManager

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	v1 "k8s.io/api/core/v1"
	"log"
	"path"
)

type VDevice struct {
	PodUID string
	vDeviceType string
	connection *utils.Messenger
	vlmMgr *VolumeManager
}

func (vd *VDevice) CheckStatus() v1.PodPhase {
	phase := utils.CheckPodStatus(vd.PodUID)
	return phase
}

func (vd *VDevice) Serve() {
	vd.connection.Actions["Register"] = vd.Register
	go vd.connection.Serve()
}

//func (vd *VDevice) monitor() {}

func (vd *VDevice) Register(PodUID string, containerID string, registerInfo string) string {
	log.Printf("Received register request from pod %s, container: %s", PodUID, containerID)
	vd.vlmMgr.UpdateInfo(containerID, PodUID)
	vd.PodUID = PodUID
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
		case v1.PodRunning:
			continue
		case v1.PodPending:
			log.Fatalf("Container Still pending after serving, please check\n")
		case v1.PodSucceeded:
			vdevice.vlmMgr.ReleaseConfig()
		default:
			log.Println("Unexpected status")
		}
	}
}



