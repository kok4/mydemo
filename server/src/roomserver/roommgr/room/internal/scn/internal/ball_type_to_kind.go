package internal

import "usercmd"
import "roomserver/roommgr/room/internal/scn/consts"

func BallTypeToKind(btype usercmd.BallType) consts.BallKind {
	if btype == usercmd.BallType_Player {
		return consts.BallKind_Player
	} else if btype > usercmd.BallType_FoodBegin && btype < usercmd.BallType_FoodEnd {
		return consts.BallKind_Food
	} else if btype > usercmd.BallType_FeedBegin && btype < usercmd.BallType_FeedEnd {
		return consts.BallKind_Feed
	} else if btype > usercmd.BallType_SkillBegin && btype < usercmd.BallType_SkillEnd {
		return consts.BallKind_Skill
	} else {
		return consts.BallKind_None
	}
}
