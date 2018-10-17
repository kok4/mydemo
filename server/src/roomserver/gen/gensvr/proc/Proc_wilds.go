package proc

// Code generated by gen.
// 本文件是对应 wilds 的处理器实现文件。
// 框架代码由 gen 生成，具体实现需要手工填写。
// 再次生成时会合并原有实现。

import (
	"base/glog"
	"roomserver/match"
	"roomserver/roommgr"
	"roomserver/roommgr/room"
	"roomserver/sess"
	"usercmd"

	"zeus/net/server"
)

// Proc_wilds 是消息处理类(Processor).
// 必须实现 NewProc_wilds(), OnClosed() 和 MsgProc_*() 接口。
type Proc_wilds struct {
	sess       server.ISession     // 一般都需要包含session对象
	sessPlayer *sess.SessionPlayer // 网络会话对应的玩家对象

	// 进入房间才会有
	room        *room.Room
	scenePlayer IScenePlayer // 场景内执行动作的玩家对象
}

type IScenePlayer interface {
	OnCastSkill(op *usercmd.MsgCastSkill)
	OnNetMove(op *usercmd.MsgMove)
	OnNetReLife(op *usercmd.MsgRelife)
	OnRun(op *usercmd.MsgRun)
	OnSceneChat(op *usercmd.MsgSceneChat)
}

var matchMaker = match.NewMatchMaker()

func NewProc_wilds(s server.ISession) *Proc_wilds {
	glog.Infof("Connected %d %s", s.GetID(), s.RemoteAddr())

	return &Proc_wilds{
		sess:       s,
		sessPlayer: sess.NewSessionPlayer(s),
	}
}

func (p *Proc_wilds) OnClosed() {
	// 会话断开时动作...
	matchMaker.DelPlayer(p.sessPlayer)

	if p.room == nil {
		return
	}

	// 下线从房间删除
	p.room.PostToRemovePlayerById(p.sessPlayer.PlayerID())
}

// IsPlaying 返回玩家是否处于正常游戏中.
func (p *Proc_wilds) IsPlaying() bool {
	return p.room != nil && !p.room.IsClosed() && p.scenePlayer != nil
}

// CheckPlaying 检查并返回玩家是否处于正常游戏中.
func (p *Proc_wilds) CheckPlaying() bool {
	if p.IsPlaying() {
		return true
	}

	if p.SetRoom() == false {
		return false
	}
	if p.SetScenePlayer() == false {
		return false
	}
	return p.IsPlaying()
}

// SetRoom 检查并设置房间. 返回是否设置.
func (p *Proc_wilds) SetRoom() bool {
	if p.room != nil {
		return true
	}

	roomID := p.sessPlayer.RoomID()
	p.room = roommgr.GetMe().GetRoomByID(roomID)
	return p.room != nil
}

// SetScenePlayer 检查并设置scenePlayer成员. 返回是否设置.
// 需要等待房间线程返回。
func (p *Proc_wilds) SetScenePlayer() bool {
	if p.scenePlayer != nil {
		return true
	}
	if p.room == nil {
		return false
	}

	// 阻塞式获取 scenePlayer 对象
	room := p.room
	playerID := p.sessPlayer.PlayerID()
	itf := room.Call(func() interface{} {
		// 在房间协程中执行
		return IScenePlayer(room.GetScenePlayer(playerID))
	})
	p.scenePlayer = itf.(IScenePlayer)
	return p.scenePlayer != nil // scenePlayer 可能还未创建
}

func (p *Proc_wilds) MsgProc_MsgLogin(msg *usercmd.MsgLogin) {
	glog.Infof("[登录] 收到登录请求 %s, %d, %s", p.sess.RemoteAddr(), p.sess.GetID(), msg.Name)
	p.sessPlayer.SetName(msg.Name)
	matchMaker.AddPlayer(p.sessPlayer)
}

func (p *Proc_wilds) MsgProc_MsgMove(msg *usercmd.MsgMove) {
	if p.CheckPlaying() == false {
		return
	}

	p.room.PostAction(func() {
		p.scenePlayer.OnNetMove(msg)
	})
}

func (p *Proc_wilds) MsgProc_MsgRun(msg *usercmd.MsgRun) {
	if p.CheckPlaying() == false {
		return
	}

	p.room.PostAction(func() {
		p.scenePlayer.OnRun(msg)
	})
}

func (p *Proc_wilds) MsgProc_MsgRelife(msg *usercmd.MsgRelife) {
	if p.CheckPlaying() == false {
		return
	}
	p.room.PostAction(func() {
		p.scenePlayer.OnNetReLife(msg)
	})
}

func (p *Proc_wilds) MsgProc_ClientHeartBeat(msg *usercmd.ClientHeartBeat) {
	p.sess.Send(msg)
}

func (p *Proc_wilds) MsgProc_MsgSceneChat(msg *usercmd.MsgSceneChat) {
	if p.CheckPlaying() == false {
		return
	}

	p.room.PostAction(func() {
		p.scenePlayer.OnSceneChat(msg)
	})
}

func (p *Proc_wilds) MsgProc_MsgCastSkill(msg *usercmd.MsgCastSkill) {
	if p.CheckPlaying() == false {
		return
	}
	p.room.PostAction(func() {
		p.scenePlayer.OnCastSkill(msg)
	})
}
