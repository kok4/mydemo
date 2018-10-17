package types

type IMsgCreator interface {
	NewMsg(MsgID) IMsg
}
