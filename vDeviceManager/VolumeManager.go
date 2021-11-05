package vDeviceManager

import (
	"fmt"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
)

const (
	vcudalibcontainerbase = "/etc/vcuda/"
	vcudalibpath = "lib"
	libcudahostpath="/root/sdcard/share/cuda-control/libcuda.so.1"
)

type VolumeManager struct {
	allocationIdentifier string
	podUID string
	containerID string
	VCudaLibHostPath string
	VCudaLibContainerPath string
	VCudaConfigHostPath string
	VCudaConfigContainerPath string
	VCudaCoreLimit int32
}

func NewVolumeManager() *VolumeManager {
	allocationIdentifier := uuid.NewV4().String()
	VCudaLibHostPath := path.Join(vcudalibcontainerbase, allocationIdentifier)
	vlmMgr := &VolumeManager{
		allocationIdentifier:     allocationIdentifier,
		containerID:                   "",
		podUID:						   "",
		VCudaLibHostPath:         VCudaLibHostPath,
		VCudaLibContainerPath:    vcudalibcontainerbase,
		VCudaConfigHostPath:      path.Join(VCudaLibHostPath, "vcuda.config"),
		VCudaConfigContainerPath: path.Join(vcudalibcontainerbase, "vcuda.config"),
	}
	vlmMgr.prepareDirectories()
	return vlmMgr
}

func (v *VolumeManager) UpdateInfo(containerID string, podUID string ) {
	v.containerID = containerID
	v.podUID = podUID

}

func (v *VolumeManager) prepareDirectories() {
	log.Printf("Making directory for %s", v.allocationIdentifier)
	err := os.Mkdir(v.VCudaLibHostPath, 0755)
	if err != nil {
		log.Errorln("Directory already exists")
		log.Errorln(err)
	}
	libcudasrc := path.Join(utils.RootPath, "tx2-k8s-device-plugin/lib/libcuda.so.1")
	log.Printf("Trying to move files from %s to %s", libcudasrc, v.VCudaLibHostPath)
	out, err := exec.Command("cp", libcudasrc, v.VCudaLibHostPath).Output()
	out, err = exec.Command("ls", v.VCudaLibHostPath).Output()
	if err != nil {
		log.Errorln(err)
	} else {
		fmt.Printf("output is %s\n", out)

	}
		//fmt.Printf("output is %s\n", out)


	//err := os.MkdirAll(v.VCudaLibHostPath, os.ModePerm)
	//if err != nil {
	//	log.Infoln("Directory exists: ")
	//	panic(err)
	//}
	//
	//out, err := exec.Command("mkdir", v.VCudaLibHostPath).Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("output is %s\n", out)
	//
	//
	//out, err = exec.Command("ls", v.VCudaLibHostPath).Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("output is %s\n", out)
	//
	//
	//libcudasrc := path.Join(utils.RootPath, "tx2-k8s-device-plugin/lib/libcuda.so.1")
	//
	//out, err = exec.Command("cp", libcudasrc, v.VCudaLibHostPath).Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("output is %s\n", out)
	//out, err = exec.Command("ls", v.VCudaLibHostPath).Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("output is %s\n", out)



	//err = copy.Copy(libcudasrc, v.VCudaLibHostPath)
	//err = os.Rename(libcudasrc, path.Join(v.VCudaLibHostPath, "libcuda.so.1"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Trying to move files from %s to %s", libcudasrc, v.VCudaLibHostPath)
}


func (v *VolumeManager) WriteConfig() {
	//dl := []byte("")
	//os.WriteFile(v.VCudaConfigContainerPath, dl, os.ModePerm)
	//vCudaConfig := C.struct_resource_data_t
	//cPodUID := C.CString(v.podUID)
	//cContName := C.CString()
	//cFileName := C.CString(filename)
}

func (v *VolumeManager) ReleaseConfig() {
	os.Remove(v.VCudaLibHostPath)
}




