package plr

// 房间玩家协议处理 辅助类

import (
	"usercmd"
)

type ScenePlayerNetMsgHelper struct {
	selfPlayer *ScenePlayer // 玩家自身的引用
}

func (this *ScenePlayerNetMsgHelper) Init(selfPlayer *ScenePlayer) {
	this.selfPlayer = selfPlayer
}

//释放技能
func (this *ScenePlayerNetMsgHelper) OnCastSkill(op *usercmd.MsgCastSkill) {
	this.selfPlayer.CastSkill(op)
}

func (this *ScenePlayerNetMsgHelper) OnNetMove(op *usercmd.MsgMove) {
	if power, angle, face, ok := this.selfPlayer.CheckMoveMsg(float64(op.Power), float64(op.Angle), op.Face); ok {
		this.selfPlayer.Move(power, angle, face)
	}
}

func (this *ScenePlayerNetMsgHelper) OnNetReLife(op *usercmd.MsgRelife) {
	this.selfPlayer.Relife()
}

// 奔跑
func (this *ScenePlayerNetMsgHelper) OnRun(op *usercmd.MsgRun) {
	this.selfPlayer.Run(op)
}

// 聊天
func (this *ScenePlayerNetMsgHelper) OnSceneChat(op *usercmd.MsgSceneChat) {
	this.selfPlayer.SendRoundMsg(op) //直接转发
}
