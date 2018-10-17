// msgreg 包注册消息, 用于发送消息时，查找消息名对应的ID.
// 生成的代码将注册所有消息。

package msgreg

import (
	"fmt"
	"reflect"
	"sync"
	"zeus/net/internal/types"
)

// Msg type -> Msg ID
// 用于发送消息
// 允许不同消息注册为同一ID, 因为可能是属于不同的服务，接收者不同。
var msgTypeToID = &sync.Map{}

func RegMsg2ID(msg types.IMsg, msgID types.MsgID) {
	actual, loaded := msgTypeToID.LoadOrStore(reflect.TypeOf(msg), msgID)
	if !loaded {
		return
	}

	actualID := actual.(types.MsgID)
	panic(fmt.Sprintf("try to register message %s(ID=%d) to ID=%d", reflect.TypeOf(msg), actualID, msgID))
}

func GetMsgID(msg types.IMsg) types.MsgID {
	ID, ok := msgTypeToID.Load(reflect.TypeOf(msg))
	if !ok {
		return 0
	}
	return ID.(types.MsgID)
}
