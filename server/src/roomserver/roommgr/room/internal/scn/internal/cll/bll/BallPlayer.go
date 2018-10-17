package bll

// 玩家球

import (
	"base/ape"
	bmath "base/math"
	"base/util"
	"roomserver/roommgr/room/internal/scn/consts"
	"roomserver/types"
	"usercmd"
)

type BallPlayer struct {
	BallMove
	Attribute
	player IScenePlayer
}

func NewBallPlayer(player IScenePlayer, ballid uint32) *BallPlayer {
	x, y := player.GetBallScene().GetRandPos()
	ball := BallPlayer{
		BallMove: BallMove{
			BallFood: BallFood{
				id:       ballid,
				Pos:      util.Vector2{x, y},
				radius:   consts.DefaultBallSize,
				BallType: usercmd.BallType_Player,
			},
		},
		player: player,
	}
	ball.Init()
	ball.ResetRect()
	return &ball
}

func (this *BallPlayer) Init() {
	this.SetMP(consts.DefaultMaxMP)
	this.SetHpMax(consts.DefaultMaxHP)
	this.SetHP(consts.DefaultMaxHP)

	this.PhysicObj = ape.NewCircleParticle(float32(this.Pos.X), float32(this.Pos.Y), float32(this.radius))
	this.player.GetBallScene().AddPlayerPhysic(this.PhysicObj)
}

func (this *BallPlayer) GetPlayerId() types.PlayerID {
	return this.player.GetID()
}

func (this *BallPlayer) GetPlayer() IScenePlayer {
	return this.player
}

//定时器1s
func (this *BallPlayer) TimeAction(nowsec int64) {

}

func (this *BallPlayer) Move(perTime float64, frameRate float64) bool {
	// 有推力情况下， 忽略原来速度方向
	if this.HasForce() == true {
		force := this.GetForce()
		pos := this.PhysicObj.GetPostion()
		this.Pos = util.Vector2{float64(pos.X), float64(pos.Y)}
		this.PhysicObj.SetVelocity(&bmath.Vector2{float32(force.X), float32(force.Y)})
		return true
	}

	pos := this.PhysicObj.GetPostion()
	this.Pos = util.Vector2{float64(pos.X), float64(pos.Y)}

	speed := consts.DefaultBallSpeed

	powerMul := util.Clamp(this.player.GetPower(), 0, 1)

	if this.player.IsRunning() {
		speed *= consts.DefaultRunRatio
		powerMul = 1
	}

	speed *= powerMul
	this.speed = *this.angleVel.MultiMethod(speed)

	vel := this.speed
	vel.ScaleBy(frameRate) //几帧执行一次物理tick
	if 0 == this.player.GetPower() {
		this.PhysicObj.SetVelocity(&bmath.Vector2{0, 0})
	} else {
		this.PhysicObj.SetVelocity(&bmath.Vector2{float32(vel.X) / 30, float32(vel.Y) / 30})
	}

	return true
}
func (this *BallPlayer) FixMapEdge() bool {
	SceneSize := this.player.GetBallScene().SceneSize()
	halfRadius := this.radius * 0.5

	if this.Pos.X < halfRadius {
		this.Pos.X = halfRadius
		this.speed.X = -this.speed.X * 0.1
	} else if this.Pos.X > SceneSize-halfRadius {
		this.Pos.X = SceneSize - halfRadius
		this.speed.X = -this.speed.X * 0.1
	}
	if this.Pos.Y < halfRadius {
		this.Pos.Y = halfRadius
		this.speed.Y = -this.speed.Y * 0.1
	} else if this.Pos.Y > SceneSize-halfRadius {
		this.Pos.Y = SceneSize - halfRadius
		this.speed.Y = -this.speed.Y * 0.1
	}

	this.rect.X = this.Pos.X
	this.rect.Y = this.Pos.Y

	return true
}

//攻击预先判断
func (this *BallPlayer) PreTryHit(target *BallPlayer) bool {
	if !this.player.GetIsLive() || !target.player.GetIsLive() {
		return false
	}
	return true
}

func (this *BallPlayer) Hit(target *BallPlayer) (int32, bool) {
	damage := consts.DefaultAttack
	target.SetHP(target.GetHP() - int32(damage))
	if target.GetHP() <= 0 {
		target.player.KilledByPlayer(this.player)
	}
	return int32(damage), true
}

func (this *BallPlayer) Eat(food *BallFood) uint32 {
	if food.GetType() == usercmd.BallType_FoodBomb {
		this.SetAttr(AttrBombNum, 1)
		this.player.RefreshPlayer()
	} else if food.GetType() == usercmd.BallType_FoodHammer {
		this.SetAttr(AttrHammerNum, 1)
		this.player.RefreshPlayer()
	}

	player := this.player

	var addexp uint32 = 0
	if food.exp != 0 {
		addexp = uint32(food.exp)
	}
	if 0 != addexp {
		player.UpdateExp(int32(addexp))
	}
	return addexp
}

func (this *BallPlayer) PreCanEat(food *BallFood) bool {
	//是否已经有锤子或者炸弹了
	if food.GetType() == usercmd.BallType_FoodHammer {
		return this.GetAttr(AttrHammerNum) == 0
	} else if food.GetType() == usercmd.BallType_FoodBomb {
		return this.GetAttr(AttrBombNum) == 0
	}
	return true
}

func (this *BallPlayer) GetEatRange() float64 {
	r := consts.DefaultEatFoodRange
	if r == 0 {
		return this.radius
	}
	return r * this.GetSizeScale()
}

func (this *BallPlayer) isNear(target *BallFood) bool {
	distance := this.SqrMagnitudeTo(target)
	eatRange := this.GetEatRange()
	return distance <= (eatRange+target.radius)*(eatRange+target.radius)
}

func (this *BallPlayer) CanEat(food *BallFood) bool {
	if !this.PreCanEat(food) {
		return false
	}
	return this.isNear(food)
}

func (this *BallPlayer) OnDead() {

}

func (this *BallPlayer) GetSizeScale() float64 {
	return 1.0
}

func (this *BallPlayer) SetAngleVelAndNormalize(x, y float64) {
	this.angleVel.X = x
	this.angleVel.Y = y
	this.angleVel.NormalizeSelf()
}
