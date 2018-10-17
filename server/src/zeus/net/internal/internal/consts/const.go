package consts

const (
	MsgIDSize = 2

	// MsgHeadSize consist message length , compression type and message id
	// 4字节长度，2字节ID号
	MsgHeadSize = 6
)

const (
	// MaxMsgBuffer 消息最大长度
	MaxMsgBuffer = 100 * 1024
)
