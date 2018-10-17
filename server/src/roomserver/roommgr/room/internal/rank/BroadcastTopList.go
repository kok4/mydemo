package rank

import (
	"roomserver/roommgr/room/internal/scn/plr"
	"roomserver/types"
	"sort"
	"time"
	"usercmd"
)

const (
	MaxTopPlayer int = 9 //排行人数
)

// BroadcastTopList 发送排行榜.
// endTime 结束时间, players 玩家对象
func BroadcastTopList(endTime int64, players map[types.PlayerID]*plr.ScenePlayer) {
	msgTop := &usercmd.MsgTop{}
	ltime := endTime - time.Now().Unix()
	if ltime > 0 {
		msgTop.EndTime = uint32(ltime)
	} else {
		msgTop.EndTime = 0
	}
	msgTop.Players = []*usercmd.MsgPlayer{}
	tmpList := _FreeRankList{}
	for _, p := range players {
		tmpList = append(tmpList, _FreeRank{PlayerID: p.ID, Score: float64(p.GetExp())})
	}
	sort.Sort(tmpList)

	var topplayers [MaxTopPlayer]usercmd.MsgPlayer // 排行榜
	for index, v := range tmpList {
		p, ok := players[v.PlayerID]
		if !ok {
			continue
		}
		if len(msgTop.Players) < MaxTopPlayer {
			topplayers[index].Id = uint64(v.PlayerID)
			topplayers[index].Name = p.Name
			topplayers[index].Curexp = uint32(v.Score)
			msgTop.Players = append(msgTop.Players, &topplayers[index])
		}
	}
	for k, v := range tmpList {
		p, ok := players[v.PlayerID]
		if !ok {
			continue
		}
		p.Rank = uint32(k + 1)
		msgTop.Rank = p.Rank
		msgTop.KillNum = uint32(v.Score)
		p.Send(msgTop)
	}
}
