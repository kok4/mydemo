package rank

import "roomserver/types"

// 自由模式排行榜
type _FreeRank struct {
	PlayerID types.PlayerID
	Score    float64
}

type _FreeRankList []_FreeRank

func (self _FreeRankList) Len() int {
	return len(self)
}

func (self _FreeRankList) Less(i, j int) bool {
	return self[i].Score > self[j].Score
}

func (self _FreeRankList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
