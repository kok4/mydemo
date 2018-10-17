// 包 plr 处理玩家类相关功能。
package plr

// 玩家类

import (
	"base/glog"
	"base/util"
	"math"
	"roomserver/roommgr/room/internal/scn/consts"
	"roomserver/roommgr/room/internal/scn/internal/cll"
	"roomserver/roommgr/room/internal/scn/internal/cll/bll"
	"roomserver/roommgr/room/internal/scn/internal/interfaces"
	"roomserver/roommgr/room/internal/scn/plr/internal"
	"roomserver/types"
	"time"
	"usercmd"
	"zeus/net/server"
)

type ScenePlayer struct {
	MoveHelper              // 检查移动消息包的辅助类
	ScenePlayerViewHelper   // 玩家视野相关辅助类
	ScenePlayerNetMsgHelper // 房间玩家协议处理辅助类
	ScenePlayerPool         // 对象池

	ID types.PlayerID // 玩家id

	Sess   ISender // 网络会话对应的玩家。 Scene.AddPlayer() 中设置
	scn    IScene  // 所在场景
	BallId uint32  // 玩家球id（一次定义，后面不变）
	Name   string  // 玩家昵称

	SelfBall *bll.BallPlayer         // 玩家球
	KillNum  uint32                  // 击杀数量
	Rank     uint32                  // 结算排名
	IsLive   bool                    // 生死
	Skill    interfaces.ISkillPlayer // 技能信息

	isMoved    bool  // 是否移动过
	isRunning  bool  // 当前是否在奔跑
	IsActClose bool  // 主动断开socket
	deadTime   int64 // 死亡时间
	msgPool    *internal.MsgPool
}

type ISender interface {
	Send(msg server.IMsg) error
	IsClosed() bool
}

// NewISkillPlayer = skill.NewISkillPlayer
var NewISkillPlayer func(player *ScenePlayer) interfaces.ISkillPlayer

// NewISkillBall = skill.NewISkillBall
var NewISkillBall func(player *ScenePlayer, ball *bll.BallSkill) interfaces.ISkillBall

func NewScenePlayer(playerID types.PlayerID, name string, scn IScene) *ScenePlayer {
	p := &ScenePlayer{
		scn:     scn,
		ID:      playerID,
		Name:    name,
		IsLive:  true,
		msgPool: internal.NewMsgPool(),
	}

	p.Init()
	return p
}

func (s *ScenePlayer) Init() {
	glog.Info("ScenePlayer.Init, id:", s.ID)
	s.ScenePlayerPool.Init()
	s.ScenePlayerNetMsgHelper.Init(s)
	s.ScenePlayerViewHelper.Init()
	s.BallId = s.GetScene().GenBallID()
	s.Skill = NewISkillPlayer(s)
	s.SelfBall = bll.NewBallPlayer(s, s.BallId)
	s.GetScene().AddBall(s.SelfBall)
	s.SelfBall.SetHP(consts.DefaultMaxHP)
	s.SelfBall.SetMP(consts.DefaultMaxMP)

}

func (s *ScenePlayer) SendChat(str string) {
	op := &usercmd.MsgSceneChat{Id: uint64(s.ID), Msg: str}
	s.BroadCastMsg(op)
}

// 释放技能
func (s *ScenePlayer) CastSkill(op *usercmd.MsgCastSkill) {
	s.Skill.CastSkill(op.Skillid, s.Face)
}

func (s *ScenePlayer) Run(op *usercmd.MsgRun) {
	if s.isRunning {
		return
	}
	if s.Power == 0 {
		return
	}
	if s.Skill.GetCurSkillId() != 0 {
		return
	}
	s.isRunning = true
}

// 移动
func (s *ScenePlayer) Move(power, angle float64, face uint32) {
	if power != 0 {
		power = 1 // power恒为1,减少移动同步影响因素
	}
	s.Power = power
	s.Face = face
	if power != 0 {
		s.Angle = angle
	}
	if power == 0 {
		s.isRunning = false
	}
}

func (s *ScenePlayer) ClientCloseSocket() {
	if s.IsActClose == true {
		return
	}
	s.IsActClose = true
	s.IsLive = false
	cmd := &usercmd.MsgActCloseSocket{}
	s.Send(cmd)
}

//死了
func (s *ScenePlayer) Dead(killer *ScenePlayer) {
	s.RealDead(killer)
}

//clearExp 是否清理经验(机器人使用)，默认false
func (s *ScenePlayer) RealDead(killer *ScenePlayer) {
	if killer != nil {
		killer.UpdateExp(consts.DefaultBallPlayerExp)
	}

	msg := &s.msgPool.MsgDeath
	msg.MaxScore = uint32(s.GetExp())
	msg.Id = uint64(s.ID)
	if killer == nil {
		msg.KillId = 0
		msg.KillName = ""
	} else {
		killer.KillNum++
		msg.KillId = uint64(killer.ID)
		msg.KillName = killer.Name
	}
	if s.Sess == nil || s.Sess.IsClosed() {
		s.scn.BroadcastMsgExcept(msg, s.ID)
	} else {
		s.BroadCastMsg(msg)
	}
	s.OnDead()
}

func (s *ScenePlayer) OnDead() {
	s.SelfBall.OnDead()
	s.IsLive = false
	s.GetScene().RemoveBall(s.SelfBall) //移除
	s.GetScene().RemovePlayerPhysic(s.SelfBall.PhysicObj)
}

func (s *ScenePlayer) GetRelifeMsg() *usercmd.MsgS2CRelife {
	msg := &s.msgPool.MsgRelife
	msg.Name = s.Name
	msg.Frame = s.GetScene().Frame()
	msg.SnapInfo = s.GetSnapInfo()
	msg.Curmp = uint32(s.SelfBall.GetMP())
	msg.Curhp = uint32(s.SelfBall.GetHP())
	return msg
}

// 复活
func (s *ScenePlayer) Relife() {
	if true == s.IsLive {
		return
	}

	s.CleanPower()
	s.IsLive = true
	s.deadTime = 0
	s.IsActClose = false

	scene := s.GetScene()

	// 添加一个新的玩家球
	exp := s.GetExp()

	ball := bll.NewBallPlayer(s, s.BallId)
	s.SelfBall = ball
	s.SetExp(exp)

	scene.AddBall(s.SelfBall) //添加一个新的

	s.SendRoundMsg(s.GetRelifeMsg()) //通知复活

	// 清除视野大小，设置新视野
	s.UpdateView(scene)
	s.UpdateViewPlayers(scene)
	s.ResetMsg()

	// 玩家视野中的所有球，发送给自己
	newMsg := &s.msgPool.MsgSceneTCP
	newMsg.Reset()
	s.LookFeeds = make(map[uint32]*bll.BallFeed)
	addfeeds, _ := s.UpdateVeiwFeeds()
	newMsg.Adds = append(newMsg.Adds, addfeeds...)

	s.LookBallSkill = make(map[uint32]*bll.BallSkill)
	adds, _ := s.UpdateVeiwBallSkill()
	newMsg.Adds = append(newMsg.Adds, adds...)

	s.LookBallFoods = make(map[uint32]*bll.BallFood)
	addfoods, _ := s.UpdateVeiwFoods()
	newMsg.Adds = append(newMsg.Adds, addfoods...)

	newMsg.AddPlayers = append(newMsg.AddPlayers, bll.PlayerBallToMsgBall(s.SelfBall))

	for _, other := range s.Others {
		newMsg.AddPlayers = append(newMsg.AddPlayers, bll.PlayerBallToMsgBall(other.SelfBall))
	}

	s.Send(newMsg)

	s.RefreshPlayer()
}

func (s *ScenePlayer) ResetMsg() {
	s.ScenePlayerPool.ResetMsg()
	s.ScenePlayerViewHelper.ResetMsg()
	s.isMoved = false
}

func (s *ScenePlayer) SendSceneMsg() {
	if s.Sess == nil || s.Sess.IsClosed() {
		return
	}

	var (
		Eats          []*usercmd.BallEat
		Adds          []*usercmd.MsgBall
		AddPlayers    []*usercmd.MsgPlayerBall
		Moves         []*usercmd.BallMove
		Hits          []*usercmd.HitMsg
		Removes       []uint32
		RemovePlayers []uint32
	)

	//feed的添加删除消息单独处理
	addfeeds, delfeeds := s.UpdateVeiwFeeds()
	Adds = append(Adds, addfeeds...)
	Removes = append(Removes, delfeeds...)

	adds, dels := s.UpdateVeiwBallSkill()
	Adds = append(Adds, adds...)
	Removes = append(Removes, dels...)

	addfoods, delfoods := s.UpdateVeiwFoods()
	Adds = append(Adds, addfoods...)
	Removes = append(Removes, delfoods...)

	addplayers, delplayers := s.updateViewBallPlayer()
	AddPlayers = append(AddPlayers, addplayers...)
	RemovePlayers = append(RemovePlayers, delplayers...)

	Eats = append(Eats, s.ScenePlayerPool.MsgEats...)
	Hits = append(Hits, s.ScenePlayerPool.MsgHits...)

	ball := s.SelfBall
	if s.isMoved {
		ballmove := s.ScenePlayerPool.MsgBallMove
		ballmove.Id = ball.GetID()
		ballmove.X = int32(ball.Pos.X * consts.MsgPosScaleRate)
		ballmove.Y = int32(ball.Pos.Y * consts.MsgPosScaleRate)

		// angle && face
		if (s.SelfBall.HasForce() == false || s.Power == 0) && s.Face != 0 {
			ballmove.Face = uint32(s.Face)
			ballmove.Angle = 0
		} else {
			ballmove.Face = 0
			ballmove.Angle = int32(s.Angle)
		}

		ballmove.State = 0
		if s.isRunning {
			ballmove.State = 2
		}
		if skillid := s.Skill.GetCurSkillId(); skillid != 0 {
			ballmove.State = skillid
		}

		Moves = append(Moves, &ballmove)
	}

	//玩家广播
	for _, other := range s.Others {
		Eats = append(Eats, other.ScenePlayerPool.MsgEats...)
		Hits = append(Hits, other.ScenePlayerPool.MsgHits...)
		if other.isMoved {
			ball = other.SelfBall
			ballmove := other.ScenePlayerPool.MsgBallMove
			ballmove.Id = ball.GetID()
			ballmove.X = int32(ball.Pos.X * consts.MsgPosScaleRate)
			ballmove.Y = int32(ball.Pos.Y * consts.MsgPosScaleRate)

			// angle && face
			if (other.SelfBall.HasForce() == false || other.Power == 0) && other.Face != 0 {
				ballmove.Face = uint32(other.Face)
				ballmove.Angle = 0
			} else {
				ballmove.Face = 0
				ballmove.Angle = int32(other.Angle)
			}

			ballmove.State = 0
			if other.isRunning {
				ballmove.State = 2
			}
			if skillid := other.Skill.GetCurSkillId(); skillid != 0 {
				ballmove.State = skillid
			}

			if other != s {
				Moves = append(Moves, &ballmove)
			}
		}
	}

	// 玩家视野中的所有消息，发送给自己
	for _, cell := range s.LookCells {
		Moves = append(Moves, cell.MsgMoves...)
	}

	if len(Adds) != 0 || len(Removes) != 0 {
		//剔除自己
		if len(Adds) != 0 {
			for k, v := range Adds {
				if v.Id == s.SelfBall.GetID() {
					Adds = append(Adds[:k], Adds[k+1:]...)
					break
				}
			}
		}

		if len(Removes) != 0 {
			for k, v := range Removes {
				if v == s.SelfBall.GetID() {
					Removes = append(Removes[:k], Removes[k+1:]...)
					break
				}
			}
		}
	}

	if len(Eats) != 0 || len(Adds) != 0 || len(AddPlayers) != 0 || len(Hits) != 0 || len(Removes) != 0 || len(RemovePlayers) != 0 {
		msg := &s.msgPool.MsgSceneTCP
		msg.Eats = Eats
		msg.Adds = Adds
		msg.AddPlayers = AddPlayers
		msg.Hits = Hits
		msg.Removes = Removes
		msg.RemovePlayers = RemovePlayers
		s.Send(msg)
	}

	if len(Moves) == 0 {
		return
	}

	msg := &s.msgPool.MsgSceneUDP
	msg.Moves = Moves
	msg.Frame = s.GetScene().Frame()
	s.Sess.Send(msg)
}

// 检查能否被吃
func (s *ScenePlayer) CanBeEat() bool {
	if s.IsLive {
		return true
	}
	return false
}

// 发送普通消息
func (s *ScenePlayer) Send(msg server.IMsg) bool {
	if s.Sess == nil || s.Sess.IsClosed() {
		return false
	}
	s.Sess.Send(msg)
	return true
}

// 广播消息
func (s *ScenePlayer) BroadCastMsg(msg server.IMsg) bool {
	s.scn.BroadcastMsg(msg)
	return true
}

// 给周围发送消息
func (s *ScenePlayer) SendRoundMsg(msg server.IMsg) bool {
	return s.AsyncRoundMsg(msg)
}

func (s *ScenePlayer) AsyncRoundMsg(msg server.IMsg) bool {
	s.Send(msg)
	for _, player := range s.RoundPlayers {
		player.Send(msg)
	}
	return true
}

func (s *ScenePlayer) UpdateMove(perTime float64, frameRate float64) {
	if !s.IsLive {
		return
	}

	// 玩家球移动
	ball := s.SelfBall
	ball.UpdateForce(perTime)
	if ball.Move(perTime, frameRate) {
		ball.FixMapEdge() //修边
		s.isMoved = true
		ball.ResetRect()

		// 扣蓝
		if s.isRunning {
			cost := frameRate * float64(consts.FrameTimeMS) * consts.DefaultRunCostMP
			diff := ball.GetMP() - cost
			if diff <= 0 {
				s.isRunning = false
			} else {
				ball.SetMP(diff)
			}
		}
	}
}

// 场景帧驱动
func (s *ScenePlayer) Update(perTime float64, now int64, scene IScene) {
	if s.Sess == nil || s.Sess.IsClosed() {
		return
	}

	curmp := s.SelfBall.GetMP()
	curexp := s.GetExp()

	// 有在释放技能，恢复转向
	s.Skill.TryTurn(&s.Angle, &s.Face)

	// 角色朝向，每帧只算一次。避免多次计算，因此代码挪至开头
	s.SelfBall.SetAngleVelAndNormalize(
		math.Cos(math.Pi*s.Angle/180),
		-math.Sin(math.Pi*s.Angle/180))

	s.Skill.Update()

	var frameRate float64 = 2

	// 更新球
	ball := s.SelfBall

	// 玩家球移动
	s.UpdateMove(perTime, frameRate)

	s.UpdateView(scene)

	if s.IsLive {
		var rect util.Square
		rect.CopyFrom(s.GetViewRect())
		rect.SetRadius(s.SelfBall.GetEatRange())
		cells := s.GetScene().GetAreaCells(&rect)
		for _, newcell := range cells {
			newcell.EatByPlayer(ball, s)
		}
	}

	// 更新视野中的玩家
	s.UpdateViewPlayers(scene)

	if curexp != s.GetExp() || curmp != s.SelfBall.GetMP() {
		s.RefreshPlayer()
	}
}

// 增加经验
func (s *ScenePlayer) UpdateExp(addexp int32) {
	if 0 == addexp {
		return
	}
	exp := s.GetExp()
	if addexp > 0 {
		exp += uint32(addexp)
	} else {
		if exp > uint32(util.AbsInt(int(addexp))) {
			exp -= uint32(util.AbsInt(int(addexp)))
		} else {
			exp = 0
		}
	}
	s.SetExp(exp)
}

func (s *ScenePlayer) GetSnapInfo() *usercmd.MsgPlayerSnap {
	msg := &s.msgPool.MsgPlayerSnap
	msg.Snapx = float32(s.SelfBall.Pos.X)
	msg.Snapy = float32(s.SelfBall.Pos.Y)
	msg.Angle = float32(s.Angle)
	msg.Id = uint64(s.ID)
	return msg
}

// 玩家定时器 (一秒一次)
func (s *ScenePlayer) TimeAction(timenow time.Time) bool {
	if false == s.IsLive {
		return true
	}

	nowsec := timenow.Unix()
	// 定时器
	s.SelfBall.TimeAction(nowsec)
	mp := int32(s.SelfBall.GetMP())
	maxmp := consts.DefaultMaxMP
	curhp := s.SelfBall.GetHP()
	maxhp := consts.DefaultMaxHP
	addmp := consts.DefaultMpRecover
	addhp := consts.DefaultHpRecover
	if addhp <= 0 {
		addhp = 1
	}
	if addmp <= 0 {
		addmp = 1
	}

	if 0 != addmp {
		if mp+int32(addmp) > int32(maxmp) {
			s.SelfBall.SetMP(float64(maxmp))
		} else {
			s.SelfBall.SetMP(float64(mp + int32(addmp)))
		}
	}
	if uint32(curhp) < uint32(maxhp) {
		if uint32(uint32(curhp)+uint32(addhp)) > uint32(maxhp) {
			s.SelfBall.SetHP(int32(maxhp))
		} else {
			s.SelfBall.SetHP(int32(curhp) + int32(addhp))
		}
	}

	s.RefreshPlayer()
	return true
}

func (s *ScenePlayer) RefreshPlayer() {
	if s.Sess == nil || s.Sess.IsClosed() {
		return
	}
	msg := &s.msgPool.MsgRefreshPlayer
	msg.Player.Id = uint64(s.ID)
	msg.Player.Name = s.Name
	msg.Player.IsLive = s.IsLive
	msg.Player.SnapInfo = s.GetSnapInfo()
	msg.Player.Curexp = s.GetExp()
	msg.Player.BallId = s.SelfBall.GetID()
	msg.Player.Curmp = uint32(s.SelfBall.GetMP())
	msg.Player.Curhp = uint32(s.SelfBall.GetHP())
	msg.Player.BombNum = int32(s.SelfBall.GetAttr(bll.AttrBombNum))
	msg.Player.HammerNum = int32(s.SelfBall.GetAttr(bll.AttrHammerNum))

	s.Sess.Send(msg)
}

//重置摇杆力
func (s *ScenePlayer) CleanPower() {
	s.Power = 0
	s.Angle = 0
}

func (s *ScenePlayer) SetIsRunning(v bool) {
	s.isRunning = v
}

func (s *ScenePlayer) GetId() types.PlayerID {
	return s.ID
}

func (s *ScenePlayer) Frame() uint32 {
	return s.scn.Frame()
}

func (s *ScenePlayer) GetExp() uint32 {
	return uint32(s.SelfBall.GetAttr(bll.AttrExp))
}

func (s *ScenePlayer) SetExp(exp uint32) {
	s.SelfBall.SetAttr(bll.AttrExp, float64(exp))
}

func (s *ScenePlayer) GetScene() IScene {
	return s.scn
}

func (s *ScenePlayer) GetBallScene() bll.IScene {
	return s.scn.(bll.IScene)
}

func (s *ScenePlayer) FindNearBallByKind(kind consts.BallKind, dir *util.Vector2, cells []*cll.Cell, ballType uint32) (interfaces.IBall, float64) {
	return s.ScenePlayerViewHelper.FindNearBallByKind(s.SelfBall, kind, dir, cells, ballType)
}

func (s *ScenePlayer) UpdateView(scene IScene) {
	if !s.IsLive {
		return
	}
	s.ScenePlayerViewHelper.UpdateView(scene, s.SelfBall, scene.SceneSize(), s.scn.CellNumX(), s.scn.CellNumY())
}

func (s *ScenePlayer) UpdateViewPlayers(scene IScene) {
	s.ScenePlayerViewHelper.UpdateViewPlayers(scene, s.SelfBall)
}

func (s *ScenePlayer) GetID() types.PlayerID {
	return s.ID
}

// 当前摇杆力度（目前恒为0或者1，来简化同步计算）
func (s *ScenePlayer) GetPower() float64 {
	return s.Power
}

func (s *ScenePlayer) IsRunning() bool {
	return s.isRunning
}

func (s *ScenePlayer) GetIsLive() bool {
	return s.IsLive
}

func (s *ScenePlayer) KilledByPlayer(killer bll.IScenePlayer) {
	s.Dead(killer.(*ScenePlayer))
}

func (s *ScenePlayer) NewSkillBall(sb *bll.BallSkill) interfaces.ISkillBall {
	return NewISkillBall(s, sb) // skill.NewISkillBall
}

func (s *ScenePlayer) GetAngle() float64 {
	return s.Angle
}

func (s *ScenePlayer) GetFace() uint32 {
	return s.Face
}

func (s *ScenePlayer) DeadTime() int64 {
	return s.deadTime
}

func (s *ScenePlayer) SetDeadTime(deadTime int64) {
	s.deadTime = deadTime
}
