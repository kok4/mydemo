package consts

//球种类
type BallKind uint8

const (
	BallKind_None   BallKind = 0
	BallKind_Player BallKind = 1  // 玩家
	BallKind_Food   BallKind = 4  // 食物
	BallKind_Feed   BallKind = 8  // 蘑菇
	BallKind_Skill  BallKind = 16 // 技能
)
