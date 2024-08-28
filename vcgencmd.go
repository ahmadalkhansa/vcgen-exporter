package main

import (
	"encoding/binary"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// https://github.com/raspberrypi/utils/blob/master/vcgencmd/vcgencmd.c
// https://github.com/raspberrypi/linux/blob/rpi-6.6.y/include/uapi/asm-generic/ioctl.h

const (
	DeviceFileName = "/dev/vcio"
)

const (
	MajorNum                = 100
	IocNrbits               = 8
	IocTypebits             = 8
	IocSizebits             = 14
	IocDirbits              = 2
	IocNrshift              = 0
	IocTypeshift            = IocNrshift + IocNrbits
	IocSizeshift            = IocTypeshift + IocTypebits
	IocDirshift             = IocSizeshift + IocSizebits
	IocWrite        uintptr = 1
	IocRead         uintptr = 2
	MaxString               = 1024
	Mode                    = os.FileMode(0666)
	GetGencmdResult         = 0x00030080
)

func VcComm(c string) (string, error) {
	var mbyte []uint32
	initbyte := []uint32{0, 0x00000000, GetGencmdResult, MaxString, 0, 0}
	device, errd := os.OpenFile(DeviceFileName, os.O_RDWR, Mode)
	if errd != nil {
		log.Fatal("device error ", errd)
	}
	defer device.Close()
	deviceFd := device.Fd()
	mbyte = append(mbyte, initbyte...)
	mbyte, mbytePtr := GenVcRequest(mbyte, c)
	erri := MboxProperty(deviceFd, mbytePtr)
	rMessage := Resp(mbyte)
	return rMessage, erri
}

func GenVcRequest(reqSlice []uint32, c string) ([]uint32, uintptr) {
	command := []byte(c)
	if len(c)%4 != 0 {
		commCorr := make([]byte, len(c)+len(c)/4+len(c)%4)
		copy(commCorr, command)
		command = commCorr
	}
	commandInt := make([]uint32, len(command)>>2)
	for x := range commandInt {
		commandInt[x] = uint32(binary.LittleEndian.Uint32(command[x<<2 : (x+1)<<2]))
	}
	reqSlice = append(reqSlice, commandInt...)
	byteOffset := make([]uint32, (len(reqSlice)+MaxString)>>2)
	reqSlice = append(reqSlice, byteOffset...)
	reqSlice[len(byteOffset)] = 0x00000000
	reqSlice[0] = uint32(len(reqSlice) * int(unsafe.Sizeof(reqSlice[0])))
	reqSliceptr := uintptr(unsafe.Pointer(&reqSlice[0]))
	return reqSlice, reqSliceptr
}

func Resp(respSlice []uint32) string {
	var message []byte
	for i := 6; respSlice[i] != 0; i++ {
		s := respSlice[i]
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(s))
		message = append(message, b...)
	}
	return string(message)
}

func IOWR(t, nr, size uintptr) uintptr {
	return IOC(IocRead|IocWrite, t, nr, size)
}

func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << IocDirshift) |
		(t << IocTypeshift) |
		(nr << IocNrshift) |
		(size << IocSizeshift)
}

func MboxProperty(fd, mb uintptr) error {
	var Ptr *string
	var IoctlMboxProperty uintptr
	IoctlMboxProperty = IOWR(MajorNum, IocNrshift, unsafe.Sizeof(Ptr))
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, IoctlMboxProperty, mb)
	var err error
	if errno != 0 {
		err = errno
	}
	return err
}

/*
func main() {
	commands := []string{
		"measure_temp",
		"measure_temp pmic",
		"get_throttled",
		"measure_clock arm core isp emmc",
		"measure_volts",
		"measure_volts sdram_c",
		"measure_volts sdram_i",
		"measure_volts sdram_p",
		"get_config total_mem",
	}
	for i := range commands {
		v := VcComm(commands[i])
		fmt.Println(v)
	}
}
*/
