// msgcrtr 包注册消息创建函数并根据ID创建消息.
/*
注册之后才可以根据消息ID创建消息，用于接收消息，
或者查找消息名对应的ID，用于发送消息。

生成的代码将注册所有消息。
*/
package msgcrtr

import (
	"fmt"
	"reflect"
	"zeus/net/internal/types"

	assert "github.com/aurelien-rainone/assertgo"
)

// 非协程安全，必须在会话开始前注册完成。
type MsgCreator struct {
	// Msg ID -> Create function
	creators map[types.MsgID]func() types.IMsg
}

func NewMsgCreator() *MsgCreator {
	return &MsgCreator{
		creators: make(map[types.MsgID]func() types.IMsg),
	}
}

func (m *MsgCreator) RegMsgCreator(msgID types.MsgID, msgCreator func() types.IMsg) {
	assert.True(msgCreator != nil, "msgCreator is nil")
	f, ok := m.creators[msgID]
	if !ok {
		m.creators[msgID] = msgCreator
		return
	}

	panic(fmt.Sprintf("try to register message ID=%d(type=%s) again to type=%s", msgID,
		reflect.TypeOf(f()), reflect.TypeOf(msgCreator())))
}

func (m *MsgCreator) NewMsg(msgID types.MsgID) types.IMsg {
	f, ok := m.creators[msgID]
	if !ok {
		return nil
	}
	return f()
}

func (m *MsgCreator) IsEmpty() bool {
	return len(m.creators) == 0
}
