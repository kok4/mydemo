package bll

import (
	"base/ape"
	"roomserver/roommgr/room/internal/scn/internal/interfaces"
	"roomserver/types"
)

type IScene interface {
	SceneID() types.SceneID
	AddBall(ball interfaces.IBall)
	AddFeedPhysic(feed ape.IAbstractParticle)
	AddPlayerPhysic(player ape.IAbstractParticle)
	GetRandPos() (x, y float64)
	SceneSize() float64
	UpdateSkillBallCell(ball *BallSkill, oldCellID int)
}
