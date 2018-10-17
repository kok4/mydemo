// 包 cll 处理场景中的格子(Cell).
package cll

// 场景中的Cell

import (
	"base/glog"
	"base/util"
	"fmt"
	"roomserver/roommgr/room/internal/scn/consts"
	"roomserver/roommgr/room/internal/scn/internal/cll/bll"
	"roomserver/roommgr/room/internal/scn/internal/interfaces"
	"usercmd"
)

var (
	CellWidth  float64 = 5 //格子像素宽
	CellHeight float64 = 5 //格子像素高
)

type Cell struct {
	id            int                        // Cell编号
	rect          util.Square                // Cell地图中的位置、区域
	Foods         map[uint32]*bll.BallFood   // 该Cell上的食物球
	playerballs   map[uint32]*bll.BallPlayer // 该Cell上的玩家球
	Feeds         map[uint32]*bll.BallFeed   // 该Cell上的动态障碍物
	Skills        map[uint32]*bll.BallSkill  // 该Cell上的技能球
	MsgMoves      []*usercmd.BallMove        // 语义上非必须，一些优化手段
	msgMovesMap   map[uint32]int             // 语义上非必须，一些优化手段
	msgRemovesMap map[uint32]bool            // 语义上非必须，一些优化手段
	msgAddsMap    map[uint32]bool            // 语义上非必须，一些优化手段
}

func NewCell(id int) *Cell {
	cell := Cell{id: id}
	cell.Clean()
	return &cell
}

func (cell *Cell) FindNearFood(player *bll.BallPlayer, pos *util.Vector2, ballType uint32, dir *util.Vector2) (*bll.BallFood, float64) {
	var min float64 = 10000
	var minball *bll.BallFood
	for _, ball := range cell.Foods {
		dis := ball.Pos.SqrMagnitudeTo(pos)
		if ballType == uint32(ball.BallType) || ballType == 0 {
			if dir != nil && util.IsSameDir(dir, ball.GetPosV(), player.GetPosV()) == false {
				continue
			}
			if player.PreCanEat(ball) {
				if minball == nil || dis < min {
					min = dis
					minball = ball
				}
			}
		}
	}
	return minball, min
}

func (cell *Cell) FindNearFeed(player *bll.BallPlayer, pos *util.Vector2, dir *util.Vector2) (*bll.BallFeed, float64) {
	var min float64
	var minball *bll.BallFeed
	for _, ball := range cell.Feeds {
		if dir != nil && util.IsSameDir(dir, ball.GetPosV(), player.GetPosV()) == false {
			continue
		}
		if player.PreCanEat(&ball.BallFood) {
			dis := ball.Pos.SqrMagnitudeTo(pos)
			if minball == nil || dis < min {
				min = dis
				minball = ball
			}
		}
	}
	return minball, min
}

func (cell *Cell) FindNearSkill(player *bll.BallPlayer, pos *util.Vector2, ballType uint32, dir *util.Vector2) (*bll.BallSkill, float64) {
	var min float64
	var minball *bll.BallSkill
	for _, ball := range cell.Skills {
		if (ballType == uint32(ball.BallType) || ballType == 0) && player.PreCanEat(&ball.BallFood) {
			if dir != nil && util.IsSameDir(dir, ball.GetPosV(), player.GetPosV()) == false {
				continue
			}
			dis := ball.Pos.SqrMagnitudeTo(pos)
			if minball == nil || dis < min {
				min = dis
				minball = ball
			}
		}
	}
	return minball, min
}

func (cell *Cell) FindNearBallByKind(player *bll.BallPlayer, pos *util.Vector2, kind consts.BallKind, dir *util.Vector2, ballType uint32) (interfaces.IBall, float64) {
	if kind == consts.BallKind_Food {
		if ball, dis := cell.FindNearFood(player, pos, ballType, dir); ball != nil {
			return ball, dis
		}
	} else if kind == consts.BallKind_Feed {
		if ball, dis := cell.FindNearFeed(player, pos, dir); ball != nil {
			return ball, dis
		}
	} else if kind == consts.BallKind_Skill {
		if ball, dis := cell.FindNearSkill(player, pos, ballType, dir); ball != nil {
			return ball, dis
		}
	}
	return nil, 10000
}

func (cell *Cell) Clean() {
	cell.Foods = make(map[uint32]*bll.BallFood)
	cell.playerballs = make(map[uint32]*bll.BallPlayer)
	cell.Feeds = make(map[uint32]*bll.BallFeed)
	cell.Skills = make(map[uint32]*bll.BallSkill)
	cell.ResetMsg()
}

func (cell *Cell) ResetMsg() {
	cell.MsgMoves = cell.MsgMoves[:0]
	cell.msgAddsMap = make(map[uint32]bool)
	cell.msgRemovesMap = make(map[uint32]bool)
	cell.msgMovesMap = make(map[uint32]int)
}

func (cell *Cell) AddMsgMove(ball interfaces.IBall) {
	x, y := ball.GetPos()
	if msgIndex, ok := cell.msgMovesMap[ball.GetID()]; ok {
		msg := cell.MsgMoves[msgIndex]
		msg.X = int32(x * consts.MsgPosScaleRate)
		msg.Y = int32(y * consts.MsgPosScaleRate)
	} else {
		cell.MsgMoves = append(cell.MsgMoves,
			&usercmd.BallMove{
				Id: ball.GetID(),
				X:  int32(x * consts.MsgPosScaleRate),
				Y:  int32(y * consts.MsgPosScaleRate),
			})
		cell.msgMovesMap[ball.GetID()] = len(cell.MsgMoves) - 1
	}
}

//添加球球
func (cell *Cell) Add(ball interfaces.IBall) {

	btype := ball.GetType()
	if btype == usercmd.BallType_Player {
		if _, ok := cell.playerballs[ball.GetID()]; !ok {
			newBall := ball.(*bll.BallPlayer)
			cell.playerballs[ball.GetID()] = newBall

			//如果已经先删除过，再添加，和之前删除抵消，不再添加
			if _, ok := cell.msgRemovesMap[ball.GetID()]; ok {
				delete(cell.msgRemovesMap, ball.GetID())
			} else {
				cell.msgAddsMap[ball.GetID()] = true
			}
		}
	} else if btype > usercmd.BallType_FoodBegin && btype < usercmd.BallType_FoodEnd {
		if _, ok := cell.Foods[ball.GetID()]; !ok {
			newBall := ball.(*bll.BallFood)
			cell.Foods[ball.GetID()] = newBall
		}
	} else if btype > usercmd.BallType_FeedBegin && btype < usercmd.BallType_FeedEnd {
		if _, ok := cell.Feeds[ball.GetID()]; !ok {
			newBall := ball.(*bll.BallFeed)
			cell.Feeds[ball.GetID()] = newBall
		}
	} else if btype > usercmd.BallType_SkillBegin && btype < usercmd.BallType_SkillEnd {
		if _, ok := cell.Skills[ball.GetID()]; !ok {
			newBall := ball.(*bll.BallSkill)
			cell.Skills[ball.GetID()] = newBall
		}
	} else {
		glog.Error("cell.Add,Fail,unknow type: ", ball.GetType(), "  tid:", ball.GetTypeId())
	}
}

//移除球球
func (cell *Cell) Remove(id uint32, typ usercmd.BallType) {
	btype := typ
	if btype == usercmd.BallType_Player {
		if _, ok := cell.playerballs[id]; ok {
			delete(cell.playerballs, id)
			//玩家的球，如果已经添加，就把添加消息删除
			if _, mok := cell.msgAddsMap[id]; mok {
				delete(cell.msgAddsMap, id)
			} else {
				cell.msgRemovesMap[id] = true
			}
		}
	} else if btype > usercmd.BallType_FoodBegin && btype < usercmd.BallType_FoodEnd {
		if _, ok := cell.Foods[id]; ok {
			delete(cell.Foods, id)
		}
	} else if btype > usercmd.BallType_FeedBegin && btype < usercmd.BallType_FeedEnd {
		if _, ok := cell.Feeds[id]; ok {
			delete(cell.Feeds, id)
		}
	} else if btype > usercmd.BallType_SkillBegin && btype < usercmd.BallType_SkillEnd {
		if _, ok := cell.Skills[id]; ok {
			delete(cell.Skills, id)
		}
	} else {
		glog.Error("[格子] 删除未知类型 ", id, ",", typ)
	}
}

//寻找球球
func (cell *Cell) Find(id uint32, typ usercmd.BallType) (interfaces.IBall, bool) {
	btype := typ
	if btype == usercmd.BallType_Player {
		ball, ok := cell.playerballs[id]
		return ball, ok
	} else if btype > usercmd.BallType_FoodBegin && btype < usercmd.BallType_FoodEnd {
		ball, ok := cell.Foods[id]
		return ball, ok
	} else if btype > usercmd.BallType_FeedBegin && btype < usercmd.BallType_FeedEnd {
		ball, ok := cell.Feeds[id]
		return ball, ok
	} else if btype > usercmd.BallType_SkillBegin && btype < usercmd.BallType_SkillEnd {
		ball, ok := cell.Skills[id]
		return ball, ok
	} else {
		return nil, false
	}
}

//全局寻找球球
func (cell *Cell) NoTypeFind(id uint32) (interfaces.IBall, bool) {
	if ball, ok := cell.playerballs[id]; ok {
		return ball, ok
	} else if ball, ok := cell.Foods[id]; ok {
		return ball, ok
	} else if ball, ok := cell.Feeds[id]; ok {
		return ball, ok
	} else if ball, ok := cell.Skills[id]; ok {
		return ball, ok
	}
	return nil, false
}

func (cell *Cell) EatByPlayer(playerBall *bll.BallPlayer, player IScenePlayer) bool {
	isEat := false
	for _, food := range cell.Foods {
		if playerBall.CanEat(food) {
			playerBall.Eat(food)
			cell.OnFoodRemoved(food)
			cell.Remove(food.GetID(), food.GetType())
			player.AddEatMsg(playerBall.GetID(), food.GetID())
			isEat = true
		}
	}
	return isEat
}

//移除食物事件
func (cell *Cell) OnFoodRemoved(food *bll.BallFood) {
	if food.GetBirthPoint() != nil {
		food.GetBirthPoint().OnChildRemove(food)
	}
}

//渲染
func (cell *Cell) Render(scene bll.IScene, per float64, now int64) {
	var delSkills []*bll.BallSkill
	for _, ball := range cell.Skills {
		if ball.Skill.IsFinish() {
			delSkills = append(delSkills, ball)
		}
	}
	for _, skill := range delSkills {
		cell.Remove(skill.GetID(), skill.GetType())
	}
	for _, ball := range cell.Skills {
		// 检查移动是否出格子
		if ball.Move(per, scene) {
			cell.AddMsgMove(ball)
			// 如果技能球已移动新的格子，则更新，删除旧格中的球，添加到新格。
			scene.UpdateSkillBallCell(ball, cell.id)
			//			x, y := ball.GetPos()
			//			newCell, ok := scene.GetCell(x, y)
			//			if ok && newCell.id != cell.id {
			//				cell.Remove(ball.GetID(), ball.GetType())
			//				newCell.Add(ball)
			//			}
			ball.ResetRect()
		}
		ball.Skill.Update()
	}
}

func (cell *Cell) string() string {
	return fmt.Sprintf("cell:%d%v", cell.id, cell.rect)
}

func (cell *Cell) ID() int {
	return cell.id
}
