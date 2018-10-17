// 包 scn 有场景类。
package scn

// 场景类

import (
	"math"
	"math/rand"
	"time"
	"zeus/net/server"

	"base/ape"
	"base/glog"
	"base/util"
	"roomserver/conf"
	"roomserver/roommgr/room/internal/scn/consts"
	"roomserver/roommgr/room/internal/scn/internal"
	"roomserver/roommgr/room/internal/scn/internal/birth"
	"roomserver/roommgr/room/internal/scn/internal/cll"
	"roomserver/roommgr/room/internal/scn/internal/cll/bll"
	"roomserver/roommgr/room/internal/scn/internal/interfaces"
	"roomserver/roommgr/room/internal/scn/internal/physic"
	"roomserver/roommgr/room/internal/scn/plr"
	"roomserver/types"
	"usercmd"

	"go.uber.org/atomic"
)

type Scene struct {
	id   types.SceneID
	room _IRoom // 所在房间

	genBallID   uint32             // 用于生成唯一BallID
	birthPoints *birth.BirthPoints // 出生点
	mapConfig   *conf.MapConfig    // 地图配置

	sceneSize float64     // 地图大小（长、宽相等）
	cellNumX  int         // 格子最大X
	cellNumY  int         // 格子最大Y
	cells     []*cll.Cell // 所有格子

	Players     map[types.PlayerID]*plr.ScenePlayer // 玩家对象
	playerCount atomic.Uint32                       // 玩家个数，= len(Players)
	scenePhysic *physic.ScenePhysic                 // 场景物理
}

func New(sceneID types.SceneID) *Scene {
	return &Scene{
		id:          sceneID,
		birthPoints: &birth.BirthPoints{},
	}
}

//场景初始化
func (s *Scene) Init(room _IRoom) {
	glog.Info("Scene Init")
	s.room = room

	s.mapConfig = conf.GetMapConfigById(s.SceneID())
	s.scenePhysic = physic.NewScenePhysic()
	s.Players = make(map[types.PlayerID]*plr.ScenePlayer)
	s.loadMap()
	for i := 0; i < s.cellNumX*s.cellNumY; i++ {
		s.cells = append(s.cells, cll.NewCell(i))
	}
	s.birthPoints.CreateAllBirthPoint(s)
}

func (s *Scene) loadMap() {
	s.sceneSize = s.mapConfig.Size
	s.cellNumX = int(math.Ceil(s.sceneSize / cll.CellWidth))
	s.cellNumY = int(math.Ceil(s.sceneSize / cll.CellHeight))
	s.scenePhysic.CreateBoard(float32(s.mapConfig.Size))
	for _, v := range s.mapConfig.Nodes {
		internal.LoadMapObjectByConfig(v, s.scenePhysic)
	}
}

//5帧更新
func (s *Scene) render5() {
	s.sendRoomMsg()
}

//时间片渲染
func (s *Scene) Render() {
	now := time.Now()
	nowNano := now.UnixNano()
	var d float64 = consts.FrameTime

	frame := s.Frame()
	if frame%2 == 0 {
		s.scenePhysic.Tick()
	}
	for _, player := range s.Players {
		player.Update(d, nowNano, s)
	}
	for _, cell := range s.cells {
		cell.Render(s, d, nowNano)
	}
	s.birthPoints.RefreshBirthPoint(d, s)

	if frame%consts.FrameCountBy100MS == 0 {
		//100ms, 5帧更新
		s.render5()
	}
}

// 发送消息, 100ms发送一次
func (s *Scene) sendRoomMsg() {
	for _, player := range s.Players {
		player.SendSceneMsg()
	}

	for _, cell := range s.cells {
		cell.ResetMsg()
	}

	for _, player := range s.Players {
		player.ResetMsg()
	}
}

//添加球到场景
func (s *Scene) AddBall(ball interfaces.IBall) {
	x, y := ball.GetPos()
	cell, ok := s.GetCell(x, y)
	if ok {
		cell.Add(ball)
	}
}

//删除球球
func (s *Scene) RemoveBall(ball interfaces.IBall) {
	if nil == ball {
		return
	}
	for _, cell := range s.cells {
		cell.Remove(ball.GetID(), ball.GetType())
	}
}

//获取区域内的所有格子
func (s *Scene) GetAreaCells(sqr *util.Square) (cells []*cll.Cell) {
	minX := int(math.Max(math.Floor(sqr.Left/cll.CellWidth), 0))
	maxX := int(math.Min(math.Floor(sqr.Right/cll.CellWidth), float64(s.cellNumX-1)))
	minY := int(math.Max(math.Floor(sqr.Bottom/cll.CellHeight), 0))
	maxY := int(math.Min(math.Floor(sqr.Top/cll.CellHeight), float64(s.cellNumY-1)))
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			cells = append(cells, s.cells[i*s.cellNumX+j])
		}
	}
	return
}

//根据坐标获取格子
func (s *Scene) GetCell(px, py float64) (*cll.Cell, bool) {
	idxX := int(math.Max(math.Floor(px/cll.CellWidth), 0))
	idxY := int(math.Max(math.Floor(py/cll.CellHeight), 0))
	//glog.Info("[房间] cells:", px, "---", py, ", " , "-", idxX, "-", idxY)
	if idxX < s.cellNumX && idxY < s.cellNumY {
		return s.cells[idxY*s.cellNumX+idxX], true
	}
	return nil, false
}

//获取场景玩家
func (s *Scene) GetPlayer(playerID types.PlayerID) *plr.ScenePlayer {
	player, ok := s.Players[playerID]
	if !ok {
		return nil
	}
	return player
}

//AddPlayer 添加玩家到场景玩家
func (s *Scene) AddPlayer(player ISessionPlayer) {
	playerID := player.PlayerID()
	_, ok := s.Players[playerID]
	if ok {
		return // 已存在
	}

	scenePlayer := s.NewScenePlayer(playerID, player.Name())
	s.Players[playerID] = scenePlayer
	s.playerCount.Store(uint32(len(s.Players)))

	scenePlayer.SetExp(0)
	scenePlayer.Sess = player

	// 把玩家发给其它人
	othermsg := &usercmd.MsgAddPlayer{}
	othermsg.Player = &usercmd.MsgPlayer{}
	othermsg.Player.Id = uint64(scenePlayer.ID)
	othermsg.Player.Name = scenePlayer.Name
	othermsg.Player.IsLive = scenePlayer.IsLive
	othermsg.Player.SnapInfo = scenePlayer.GetSnapInfo()
	othermsg.Player.BallId = scenePlayer.SelfBall.GetID()
	othermsg.Player.Curmp = uint32(scenePlayer.SelfBall.GetMP())
	othermsg.Player.Curhp = uint32(scenePlayer.SelfBall.GetHP())

	scenePlayer.UpdateView(s)
	scenePlayer.UpdateViewPlayers(s)
	scenePlayer.ResetMsg()

	// 发送MsgTop消息给玩家(主要是更新EndTime)
	msgTop := &usercmd.MsgTop{}
	ltime := s.room.EndTime() - time.Now().Unix()
	if ltime > 0 {
		msgTop.EndTime = uint32(ltime)
	} else {
		msgTop.EndTime = 0
	}
	scenePlayer.Sess.Send(msgTop)

	// 把当前场景的人、球都发给玩家
	var others []*usercmd.MsgPlayer
	var balls []*usercmd.MsgBall
	var playerballs []*usercmd.MsgPlayerBall
	for _, player := range s.Players {
		others = append(others, &usercmd.MsgPlayer{
			Id:     uint64(player.ID),
			BallId: player.SelfBall.GetID(),
			Name:   player.Name,
			IsLive: player.IsLive,

			SnapInfo: player.GetSnapInfo(),
			Curhp:    uint32(player.SelfBall.GetHP()),
			Curmp:    uint32(player.SelfBall.GetMP()),
			Curexp:   player.GetExp(),
		})
	}

	// 玩家视野中的所有球，发送给自己
	cells := scenePlayer.LookCells

	scenePlayer.LookFeeds = make(map[uint32]*bll.BallFeed)
	addfeeds, _ := scenePlayer.UpdateVeiwFeeds()
	balls = append(balls, addfeeds...)

	scenePlayer.LookBallSkill = make(map[uint32]*bll.BallSkill)
	adds, _ := scenePlayer.UpdateVeiwBallSkill()
	balls = append(balls, adds...)

	scenePlayer.LookBallFoods = make(map[uint32]*bll.BallFood)
	addfoods, _ := scenePlayer.UpdateVeiwFoods()
	balls = append(balls, addfoods...)

	//自己
	playerballs = append(playerballs, bll.PlayerBallToMsgBall(scenePlayer.SelfBall))
	//周围玩家
	for _, other := range scenePlayer.Others {
		if true == other.IsLive {
			playerballs = append(playerballs, bll.PlayerBallToMsgBall(other.SelfBall))
		}
	}

	msg := &usercmd.MsgLoginResult{}
	msg.Id = uint64(scenePlayer.ID)
	msg.BallId = scenePlayer.SelfBall.GetID()
	msg.Name = scenePlayer.Name
	msg.Ok = true
	msg.Frame = s.Frame()
	msg.Balls = balls
	msg.Playerballs = playerballs
	msg.Others = others
	msg.LeftTime = uint32(s.room.EndTime() - time.Now().Unix())

	scenePlayer.Sess.Send(msg)
	glog.Info("[登录] 添加玩家成功addplayer [", s.room.ID(), ",", scenePlayer.Name, "],", scenePlayer.ID, ",ballId:", msg.BallId, ",view:", scenePlayer.GetViewRect(), ",cell:", len(cells),
		",otplayer:", len(others), "so", len(scenePlayer.Others), ",ball:", len(balls))

	s.BroadcastMsg(othermsg)
}

// 删除玩家
func (s *Scene) RemovePlayer(playerId types.PlayerID) bool {
	player, ok := s.Players[playerId]
	if !ok {
		return false
	}
	delete(s.Players, playerId)
	s.playerCount.Store(uint32(len(s.Players)))

	oldstatus := player.IsLive

	s.RemoveBall(player.SelfBall)
	s.scenePhysic.RemovePlayer(player.SelfBall.PhysicObj)

	player.IsLive = false
	player.SetDeadTime(time.Now().Unix())
	player.Dead(nil)
	player.Sess = nil

	glog.Info("[注销] 删除玩家成功 [", s.room.ID(), "],", player.ID, " players:", len(s.Players), ";", oldstatus, ",exp:", player.GetExp())
	return true
}

func (s *Scene) SceneID() types.SceneID {
	return s.id
}

func (s *Scene) RemoveFeed(feed *bll.BallFeed) {
	if feed.PhysicObj != nil {
		s.scenePhysic.RemoveFeed(feed.PhysicObj)
	}
}

func (s *Scene) AddFeedPhysic(feed ape.IAbstractParticle) {
	s.scenePhysic.AddFeed(feed)
}

func (s *Scene) AddPlayerPhysic(player ape.IAbstractParticle) {
	s.scenePhysic.AddPlayer(player)
}

func (s *Scene) RemovePlayerPhysic(player ape.IAbstractParticle) {
	s.scenePhysic.RemovePlayer(player)
}

// 地图大小（长、宽相等）
func (s *Scene) SceneSize() float64 {
	return s.sceneSize
}

func (s *Scene) UpdateSkillBallCell(ball *bll.BallSkill, oldCellID int) {
	x, y := ball.GetPos()
	newCell, ok := s.GetCell(x, y)
	if !ok || newCell.ID() == oldCellID {
		return
	}
	oldCell := s.cells[oldCellID]
	if oldCell.ID() != oldCellID {
		panic("从Cell ID获取相应Cell算法有误")
	}

	oldCell.Remove(ball.GetID(), ball.GetType())
	newCell.Add(ball)
}

func (s *Scene) Frame() uint32 {
	return s.room.Frame()
}

// 格子最大X
func (s *Scene) CellNumX() int {
	return s.cellNumX
}

// 格子最大Y
func (s *Scene) CellNumY() int {
	return s.cellNumY
}

func (s *Scene) GetPlayers() map[types.PlayerID]*plr.ScenePlayer {
	return s.Players
}

// 协程安全
func (s *Scene) GetPlayerCount() uint32 {
	return s.playerCount.Load()
}

// 广播 msg
func (s *Scene) BroadcastMsg(msg server.IMsg) {
	for _, c := range s.Players {
		c.Send(msg)
	}
}

//广播(剔除特定ID)
func (s *Scene) BroadcastMsgExcept(msg server.IMsg, uid types.PlayerID) {
	for _, c := range s.Players {
		if c.ID == uid {
			continue
		}
		c.Send(msg)
	}
}

func (s *Scene) NewScenePlayer(playerID types.PlayerID, name string) *plr.ScenePlayer {
	return plr.NewScenePlayer(playerID, name, s)
}

func (s *Scene) GetRandPos() (x, y float64) {
	return s.sceneSize * rand.Float64(), s.sceneSize * rand.Float64()
}

func (s *Scene) GenBallID() uint32 {
	s.genBallID++
	return s.genBallID
}
