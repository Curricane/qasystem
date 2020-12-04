package id_gen

import (
	"fmt"
	"github.com/sony/sonyflake"
)

var (
	sonyFlake *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

func Init(machineId uint16) (err error) {
	sonyMachineID = machineId
	settings := sonyflake.Settings{}
	settings.MachineID = getMachineID // 通过回调函数的方式获取MachineID

	sonyFlake = sonyflake.NewSonyflake(settings)

	return
}

func GetId() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sonyFlake not inited")
		return
	}
	id, err = sonyFlake.NextID()
	return
}
