package skill

// 几种选取目标的方式。如用于释放技能时

import (
	b3core "base/behavior3go/core"
	_ "base/glog"
	"base/util"
	"roomserver/roommgr/room/internal/scn/consts"
	"roomserver/roommgr/room/internal/scn/internal/interfaces"
	"roomserver/roommgr/room/internal/scn/internal/cll/bll"
	"roomserver/roommgr/room/internal/scn/plr"
)

// 获取朝向上最近的目标
func FindNearTarget(tick *b3core.Tick, player *plr.ScenePlayer) (interfaces.IBall, consts.BallKind) {
	angleVel := GetPlayerDir(tick, player)

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(GetAttackRange(tick, player))
	cells := player.GetScene().GetAreaCells(&rect)

	minball_feed, min_feed := player.FindNearBallByKind(consts.BallKind_Feed, angleVel, cells, 0)
	minball_player, min_player := player.FindNearBallByKind(consts.BallKind_Player, angleVel, cells, 0)
	minball_ballskill, min_ballskill := player.FindNearBallByKind(consts.BallKind_Skill, angleVel, cells, 0)
	if minball_player == nil && minball_feed == nil && minball_ballskill == nil {
		return nil, consts.BallKind_None
	}

	if min_feed <= min_player && min_feed <= min_ballskill {
		return minball_feed, consts.BallKind_Feed
	} else if min_ballskill <= min_player && min_ballskill <= min_feed {
		return minball_ballskill, consts.BallKind_Skill
	} else {
		return minball_player, consts.BallKind_Player
	}
}

// 获取朝向上所有目标
func FindTarget_SemiCircle(tick *b3core.Tick, player *plr.ScenePlayer) ([]interfaces.IBall, []consts.BallKind) {
	var balllist []interfaces.IBall
	var balltype []consts.BallKind

	dir := GetPlayerDir(tick, player)

	// player
	for _, o := range player.Others {
		if o.IsLive == false {
			continue
		}
		ball := o.SelfBall
		if util.IsSameDir(dir, ball.GetPosV(), player.SelfBall.GetPosV()) == false {
			continue
		}
		balllist = append(balllist, ball)
		balltype = append(balltype, consts.BallKind_Player)
	}

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(GetAttackRange(tick, player))
	cells := player.GetScene().GetAreaCells(&rect)

	// ballskill
	for _, cell := range cells {
		for _, ball := range cell.Skills {
			if util.IsSameDir(dir, ball.GetPosV(), player.SelfBall.GetPosV()) == false {
				continue
			}
			balllist = append(balllist, ball)
			balltype = append(balltype, consts.BallKind_Skill)
		}
	}

	// feed
	for _, cell := range cells {
		for _, ball := range cell.Feeds {
			if util.IsSameDir(dir, ball.GetPosV(), player.SelfBall.GetPosV()) == false {
				continue
			}
			balllist = append(balllist, ball)
			balltype = append(balltype, consts.BallKind_Feed)
		}
	}

	return balllist, balltype
}

// 获取所有目标
func FindTarget_Circle(tick *b3core.Tick, player *plr.ScenePlayer) ([]interfaces.IBall, []consts.BallKind) {
	var balllist []interfaces.IBall
	var balltype []consts.BallKind

	// player
	for _, o := range player.Others {
		if o.IsLive == false {
			continue
		}
		ball := o.SelfBall
		balllist = append(balllist, ball)
		balltype = append(balltype, consts.BallKind_Player)
	}

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(GetAttackRange(tick, player))
	cells := player.GetScene().GetAreaCells(&rect)

	// ballskill
	for _, cell := range cells {
		for _, ball := range cell.Skills {
			balllist = append(balllist, ball)
			balltype = append(balltype, consts.BallKind_Skill)
		}
	}

	// feed
	for _, cell := range cells {
		for _, ball := range cell.Feeds {
			balllist = append(balllist, ball)
			balltype = append(balltype, consts.BallKind_Feed)
		}
	}
	return balllist, balltype
}

// 获取玩家朝向
func GetPlayerDir(tick *b3core.Tick, player *plr.ScenePlayer) *util.Vector2 {
	angleVel := &util.Vector2{}
	usedefault := true
	targetId := tick.Blackboard.GetUInt32("skillTargetId", "", "")
	if 0 != targetId {
		tball := player.FindViewPlayer(targetId)
		if tball != nil {
			x, y := tball.GetPos()
			angleVel.X = x - player.SelfBall.GetPosV().X
			angleVel.Y = y - player.SelfBall.GetPosV().Y
			usedefault = false
		}
	}
	if usedefault {
		angleVel.X = player.SelfBall.GetAngleVel().X
		angleVel.Y = player.SelfBall.GetAngleVel().Y
	}
	return angleVel
}

// 攻击范围
func GetAttackRange(tick *b3core.Tick, player *plr.ScenePlayer) float64 {
	attackRange := tick.Blackboard.Get("attackRange", "", "")
	if attackRange != nil {
		r := attackRange.(float64)
		if r >= 0 {
			return r * player.SelfBall.GetSizeScale()
		}
	}
	return player.SelfBall.GetEatRange()
}

// 是否可以攻击
func IsCanAttack(tick *b3core.Tick, player *plr.ScenePlayer, target interfaces.IBall) bool {
	distance := player.SelfBall.SqrMagnitudeTo(target)
	eatRange := GetAttackRange(tick, player)
	return distance <= (eatRange+target.GetRect().Radius)*(eatRange+target.GetRect().Radius)
}

func IsCanAttackPlayer(tick *b3core.Tick, player *plr.ScenePlayer, target *bll.BallPlayer) bool {
	if player.SelfBall.PreTryHit(target) == false {
		return false
	}
	distance := player.SelfBall.SqrMagnitudeTo(target)
	eatRange := GetAttackRange(tick, player)
	return distance <= (eatRange+target.GetRect().Radius)*(eatRange+target.GetRect().Radius)
}
