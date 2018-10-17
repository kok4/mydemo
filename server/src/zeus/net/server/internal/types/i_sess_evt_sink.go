package types

type ISessEvtSink interface {
	OnConnected(ISession)
	OnClosed(ISession)
}
