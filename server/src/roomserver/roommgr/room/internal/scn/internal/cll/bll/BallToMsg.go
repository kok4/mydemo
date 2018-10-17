package bll

// 球信息 到 球网络信息

import (
	"roomserver/roommgr/room/internal/scn/consts"
	"usercmd"
)

func FoodToMsgBall(ball *BallFood) *usercmd.MsgBall {
	return &usercmd.MsgBall{
		Id:   ball.id,
		Type: int32(ball.typeID),
		X:    int32(ball.Pos.X * consts.MsgPosScaleRate),
		Y:    int32(ball.Pos.Y * consts.MsgPosScaleRate),
	}
}

func FeedToMsgBall(ball *BallFeed) *usercmd.MsgBall {
	cmd := &usercmd.MsgBall{
		Id:   ball.id,
		Type: int32(ball.typeID),
		X:    int32(ball.Pos.X * consts.MsgPosScaleRate),
		Y:    int32(ball.Pos.Y * consts.MsgPosScaleRate),
	}
	return cmd
}

func SkillToMsgBall(ball *BallSkill) *usercmd.MsgBall {
	cmd := &usercmd.MsgBall{
		Id:   ball.id,
		Type: int32(ball.BallType),
		X:    int32(ball.Pos.X * consts.MsgPosScaleRate),
		Y:    int32(ball.Pos.Y * consts.MsgPosScaleRate),
	}
	return cmd
}

func PlayerBallToMsgBall(ball *BallPlayer) *usercmd.MsgPlayerBall {
	cmd := &usercmd.MsgPlayerBall{
		Id:    ball.id,
		Hp:    uint32(ball.GetHP()),
		Mp:    uint32(ball.GetMP()),
		X:     int32(ball.Pos.X * consts.MsgPosScaleRate),
		Y:     int32(ball.Pos.Y * consts.MsgPosScaleRate),
		Angle: int32(ball.player.GetAngle()),
		Face:  uint32(ball.player.GetFace()),
	}

	return cmd
}
