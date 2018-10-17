package scn

import (
	"base/glog"
	"roomserver/roommgr/room/internal/scn/internal/skill"
	"roomserver/roommgr/room/internal/scn/plr"
)

func init() {
	plr.NewISkillPlayer = skill.NewISkillPlayer
	plr.NewISkillBall = skill.NewISkillBall
}

// LoadSkillBevTree 初始化加载技能行为树.
// 需要在配置加载之后执行，所以不能在 init() 中。
func LoadSkillBevTree() bool {
	if skill.LoadSkillBevTree() {
		return true
	}
	glog.Error("[启动]LoadSkillBevTree fail! ")
	return false
}
