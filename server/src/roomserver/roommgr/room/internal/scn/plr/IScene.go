package plr

import (
	"base/ape"
	"base/util"
	"roomserver/roommgr/room/internal/scn/internal/cll"
	"roomserver/roommgr/room/internal/scn/internal/cll/bll"
	"roomserver/roommgr/room/internal/scn/internal/interfaces"
	"roomserver/types"
	"zeus/net/server"
)

type IScene interface {
	GenBallID() uint32 // 生成一个场景内唯一的BallID
	AddBall(ball interfaces.IBall)
	RemoveBall(ball interfaces.IBall)
	RemovePlayerPhysic(player ape.IAbstractParticle)
	Frame() uint32
	GetAreaCells(s *util.Square) (cells []*cll.Cell)
	GetPlayers() map[types.PlayerID]*ScenePlayer
	SceneSize() float64
	SceneID() types.SceneID
	GetCell(px, py float64) (*cll.Cell, bool)
	RemoveFeed(feed *bll.BallFeed)

	CellNumX() int
	CellNumY() int

	BroadcastMsg(msg server.IMsg)
	BroadcastMsgExcept(msg server.IMsg, uid types.PlayerID)
}
