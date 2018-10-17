package interfaces

// 球的基本接口

import (
	"base/util"
	"usercmd"
)

type IBall interface {
	GetID() uint32
	SetID(uint32)
	GetTypeId() uint16
	GetType() usercmd.BallType
	GetPos() (float64, float64)
	SetPos(float64, float64)
	SetBirthPoint(_birthPoint IBirthPoint)
	GetRect() *util.Square
	OnReset()
}
