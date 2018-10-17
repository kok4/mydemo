package scn

import (
	"roomserver/roommgr/room/internal/scn/plr"
	"roomserver/types"
)

type ISessionPlayer interface {
	plr.ISender
	Name() string
	PlayerID() types.PlayerID
}
