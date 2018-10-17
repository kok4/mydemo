package room

// 房间类

import (
	"base/glog"
	"roomserver/roommgr/room/internal/rank"
	"roomserver/roommgr/room/internal/scn"
	"roomserver/roommgr/room/internal/scn/consts"
	"roomserver/roommgr/room/internal/scn/plr"
	"roomserver/types"
	"runtime/debug"
	"time"
	"usercmd"
	"zeus/net/server"

	"go.uber.org/atomic"

	assert "github.com/aurelien-rainone/assertgo"
)

type Room struct {
	id  types.RoomID // 房间id
	scn *scn.Scene   // 场景信息

	endTime int64  // 结束时间
	frame   uint32 // 当前帧数

	actC  chan func() // 玩家输入或其他动作，需要在房间协程中执行
	doneC chan bool   // 用于结束房间协程

	isClosed atomic.Bool // 是否关闭标识
}

var genRoomID types.RoomID

func Init() bool {
	return scn.LoadSkillBevTree()
}

// 创建房间
func New(sceneID types.SceneID) *Room {
	id := generateRoomID()
	glog.Infof("[房间] 创建, ID=%d", id)
	room := &Room{
		id:      id,
		scn:     scn.New(sceneID),
		endTime: time.Now().Unix() + consts.DefaultPlayTime, // 10min为一局
		actC:    make(chan func(), 1024),
		doneC:   make(chan bool),
	}
	room.scn.Init(room)

	// 开启逻辑处理协程
	go room.Loop()

	return room
}

func generateRoomID() types.RoomID {
	genRoomID++
	return genRoomID
}

// 关闭房间
func (r *Room) close() {
	if !r.isClosed.CAS(false /*old*/, true /*new*/) {
		return // CompareAndSet()失败说明已经关闭了
	}
	glog.Infof("[房间] 关闭, ID=%d, players=%d", r.id, len(r.scn.Players))

	// 更新排行榜
	r.broadcastTopList()
	// 广播房间结束
	r.broadcastEndMsg()

	close(r.doneC) // 结束房间协程
}

func (r *Room) IsClosed() bool {
	return r.isClosed.Load()
}

//主循环
func (r *Room) Loop() {
	timeTicker := time.NewTicker(time.Millisecond * consts.FrameTimeMS)

	defer func() {
		timeTicker.Stop()
		if err := recover(); err != nil {
			glog.Error("[异常] 房间线程出错 [", r.id, "] ", err, "\n", string(debug.Stack()))
		}
	}()

	for {
		select {
		case <-timeTicker.C:
			r.render()
		case act := <-r.actC:
			act()
		case <-r.doneC:
			assert.True(r.IsClosed())
			return
		}
	}
}

//添加玩家
func (r *Room) AddPlayer(player scn.ISessionPlayer) {
	// 添加到场景
	r.scn.AddPlayer(player)
	glog.Info("[登录] 进入房间成功 [", r.id, "],", player.PlayerID())
}

// 从房间里删除玩家
func (r *Room) removePlayerById(playerId types.PlayerID) {
	//退出房间处理
	r.scn.RemovePlayer(playerId)

	// 通知其它人删除玩家
	rmCmd := &usercmd.MsgRemovePlayer{
		Id: uint64(playerId),
	}
	r.BroadcastMsg(rmCmd)

	glog.Info("[房间] 删除玩家 [", r.id, "] ", playerId, ",", len(r.scn.Players))
}

// 广播结束消息
func (r *Room) broadcastEndMsg() {
	// 游戏结束，可以对 MsgEndRoom 字段做修改，发送具体游戏结果
	msg := &usercmd.MsgEndRoom{}
	r.BroadcastMsg(msg)
}

// 发送排行榜
func (r *Room) broadcastTopList() {
	rank.BroadcastTopList(r.endTime, r.scn.Players)
}

// 广播 msg
func (r *Room) BroadcastMsg(msg server.IMsg) {
	r.scn.BroadcastMsg(msg)
}

//广播(剔除特定ID)
func (r *Room) BroadcastMsgExcept(msg server.IMsg, uid types.PlayerID) {
	r.scn.BroadcastMsgExcept(msg, uid)
}

// 房间帧定时器
func (r *Room) render() {
	r.frame++
	r.scn.Render()

	//1s
	if r.frame%(consts.FrameCountBy100MS*10) != 0 {
		return
	}
	r.timeAction1s()

	if r.frame%(consts.FrameCountBy100MS*40) == 0 {
		r.broadcastTopList()
	}
}

// 房间定时器 (一秒一次)
func (r *Room) timeAction1s() {
	timeNow := time.Now()
	if r.endTime < timeNow.Unix() || len(r.scn.Players) == 0 {
		r.close()
		return
	}

	for _, player := range r.scn.Players {
		player.TimeAction(timeNow)
	}
}

// PostAction 发送动作到 ActC. 动作在房间协程中执行。
func (r *Room) PostAction(action func()) {
	r.actC <- action
}

// Call 发送函数到 ActC 并等待返回值. 函数在房间协程中执行。
func (r *Room) Call(fn func() interface{}) interface{} {
	ch := make(chan interface{}, 1) // 用于返回结果
	r.actC <- func() {
		ch <- fn()
	}
	return <-ch // 等待直到返回结果
}

func (r *Room) ID() types.RoomID {
	return r.id
}

// PostToRemovePlayerById 计划移除玩家，将在房间协程中执行动作。
func (r *Room) PostToRemovePlayerById(playerID types.PlayerID) {
	r.PostAction(func() {
		r.removePlayerById(playerID)
	})
}

func (r *Room) EndTime() int64 {
	return r.endTime
}

// 协程安全
func (r *Room) GetPlayerCount() uint32 {
	return r.scn.GetPlayerCount()
}

//获取场景玩家
func (r *Room) GetScenePlayer(id types.PlayerID) *plr.ScenePlayer {
	return r.scn.GetPlayer(id)
}

func (r *Room) Frame() uint32 {
	return r.frame
}
