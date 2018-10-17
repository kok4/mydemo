package scn

import (
	"roomserver/types"
)

type _IRoom interface {
	ID() types.RoomID
	EndTime() int64
	Frame() uint32
}
