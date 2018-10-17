package sess

import (
	"zeus/net/server"

	"roomserver/types"
)

// 会话玩家类.
type SessionPlayer struct {
	sess server.ISession

	id     types.PlayerID
	name   string
	roomID types.RoomID // 房间ID
}

// NewSessionPlayer 创建会话玩家.
// 随会话创建.
func NewSessionPlayer(sess server.ISession) *SessionPlayer {
	return &SessionPlayer{
		sess: sess,
		id:   types.PlayerID(sess.GetID()), // 以会话ID作为玩家ID
	}
}

// Send 发送消息.
func (p *SessionPlayer) Send(msg server.IMsg) error {
	p.sess.Send(msg)
	return nil
}

func (p *SessionPlayer) Name() string {
	return p.name
}

func (p *SessionPlayer) SetName(name string) {
	p.name = name
}

// PlayerID 返回玩家ID.
// 同网络会话的ID，每个连接都不同。
func (p *SessionPlayer) PlayerID() types.PlayerID {
	return p.id
}

func (p *SessionPlayer) SetRoomID(roomID types.RoomID) {
	p.roomID = roomID
}

func (p *SessionPlayer) IsClosed() bool {
	if p.sess != nil {
		return p.sess.IsClosed()
	}
	return true
}

func (p *SessionPlayer) RoomID() types.RoomID {
	return p.roomID
}
