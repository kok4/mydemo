// match 包处理匹配.
package match

import (
	"roomserver/roommgr"
	"roomserver/sess"
	"sync"
)

const DEFAULT_MATCH_NUM = 2 // 几人匹配改这里
const kSceneID = 1002       // 目前只有一个地图

type MatchMaker struct {
	mtx sync.Mutex

	players []*sess.SessionPlayer
}

func NewMatchMaker() *MatchMaker {
	return &MatchMaker{}
}

func (m *MatchMaker) AddPlayer(player *sess.SessionPlayer) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.players = append(m.players, player)
	if len(m.players) < DEFAULT_MATCH_NUM {
		return // 人还不够
	}

	// 开始游戏
	r := roommgr.GetMe().NewRoom(kSceneID)
	for _, player := range m.players {
		player := player // 必须捕获迭代变量！
		player.SetRoomID(r.ID())
		r.PostAction(func() {
			r.AddPlayer(player) // 须在房间协程中执行
		})
	}

	// 清空队列
	m.players = nil
}

func (m *MatchMaker) DelPlayer(player *sess.SessionPlayer) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for i, item := range m.players {
		if item != player {
			continue
		}
		// 删除该元素
		m.players = append(m.players[:i], m.players[i+1:]...)
		return
	}
}
