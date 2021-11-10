package vDeviceManager

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/exec"
	"path"
)
// DO NEVER DELETE THESE C LINES!!!
// THESE EMBEDDED C CODES RUN!!!

//#include <stdint.h>
//#include <sys/types.h>
//#include <sys/stat.h>
//#include <fcntl.h>
//#include <string.h>
//#include <sys/file.h>
//#include <time.h>
//#include <stdlib.h>
//#include <unistd.h>
//
//
//#ifndef FILENAME_MAX
//#define FILENAME_MAX 4096
//#endif
//typedef struct {
//  int major;
//  int minor;
//} __attribute__((packed, aligned(8))) version_t;
//
//typedef struct resource_data_t{
//char pod_uid[48];
//int limit;
//char occupied[4044];
//char container_name[FILENAME_MAX];
//uint64_t gpu_memory;
//int utilization;
//int hard_limit;
//version_t driver_version;
//int enable;
//} __attribute__((packed, aligned(8))) ;
//
//int setting_to_disk(const char* filename, struct resource_data_t* data) {
//  int fd = 0;
//  int wsize = 0;
//  int ret = 0;
//
//  fd = open(filename, O_CREAT | O_TRUNC | O_WRONLY, 00777);
//  if (fd == -1) {
//    return 1;
//  }
//
//  wsize = (int)write(fd, (void*)data, sizeof(struct resource_data_t));
//  if (wsize != sizeof(struct resource_data_t)) {
//    ret = 2;
//	goto DONE;
//  }
//
//DONE:
//  close(fd);
//
//  return ret;
//}
//
import "C"


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

		//VCudaConfigHostPath:      path.Join(VCudaLibHostPath, "vcuda.config"),
		//VCudaConfigContainerPath: path.Join(vcudalibcontainerbase, "vcuda.config"),
		VCudaConfigHostPath:      "",
		VCudaConfigContainerPath: "",
	}
	vlmMgr.prepareDirectories()
	return vlmMgr
}

func (v *VolumeManager) UpdateInfo(containerID string, podUID string ) {
	v.containerID = containerID
	v.podUID = podUID
	v.VCudaConfigHostPath = path.Join(v.VCudaLibHostPath, containerID)
	v.VCudaConfigContainerPath = path.Join(v.VCudaLibHostPath, containerID)
	err := os.Mkdir(v.VCudaConfigHostPath, 0755)
	if err != nil {
		log.Errorln("Directory already exists")
		log.Errorln(err)
	}

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
	_, err = exec.Command("cp", libcudasrc, v.VCudaLibHostPath).Output()
	_, err = exec.Command("ls", v.VCudaLibHostPath).Output()
	if err != nil {
		log.Errorln(err)
	}

}



func (v *VolumeManager) writeToDisk(cores int64, hardlimit int64,
	podUID string, containerID string , filename string) error {

	var vcudaConfig C.struct_resource_data_t
	var vcudaConfigPtr *C.struct_resource_data_t = &vcudaConfig
	vcudaConfig.limit = C.int(cores)
	cPodUID := C.CString(podUID)
	cContName := C.CString(containerID)
	cFileName := C.CString(filename)
	C.strcpy(&(vcudaConfigPtr.pod_uid[0]), (*C.char)(unsafe.Pointer(cPodUID)))
	C.strcpy(&(vcudaConfigPtr.container_name[0]), (*C.char)(unsafe.Pointer(cContName)))
	vcudaConfig.utilization = C.int(cores)
	//vcudaConfig.hard_limit = C.int(hardlimit)
	vcudaConfig.enable = 1
	if C.setting_to_disk(cFileName, vcudaConfigPtr) != 0 {
		return fmt.Errorf("can't sink config %s", filename)
	}
	return nil
}

func (v *VolumeManager) WriteConfig() error {
	os.Mkdir(v.VCudaConfigHostPath, 0755)
	configFilename := path.Join(v.VCudaConfigHostPath, "vcuda.config")

	//pod, err := utils.K8sClient.CoreV1().Pods("").Get(context.Background(), v.getPodName(v.podUID),  metav1.GetOptions{})
	pods, err := utils.K8sClient.CoreV1().Pods("").List(context.Background(),  metav1.ListOptions{})
	utils.Check(err)
	var targetPod v1.Pod
	var container v1.Container
	for _, pod := range pods.Items {
		if pod.UID == types.UID(v.podUID) {
			targetPod = pod
			container = targetPod.Spec.Containers[0]
			break
		}
	}

	coreLimit := container.Resources.Limits["csu.ac.cn/gpu"]
	coreLimitData := int64(coreLimit.Value())
	//log.Println(coreLimitData)
	//log.Println(configFilename)
	//log.Println(container.Resources.Limits["csu.ac.cn/gpu"])
	err = v.writeToDisk(coreLimitData, 1, v.podUID, v.containerID, configFilename)
	utils.Check(err)
	return nil
}

func (v *VolumeManager) ReleaseConfig() {
	log.Printf("Trying to remove %s\n", v.VCudaLibHostPath)
	_, err := exec.Command("rm", "-rf", v.VCudaLibHostPath).Output()
	if err != nil {
		log.Errorln(err)
	}
	log.Printf("%s has been removed\n", v.VCudaLibHostPath)
}




