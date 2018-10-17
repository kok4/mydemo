package server

const (
	_ = iota

	// ServerTypeClient 客户端
	ServerTypeClient
	// ServerTypeLobby 网关服
	ServerTypeLobby
	// ServerTypeMatch 匹配服
	ServerTypeMatch
)

const (
	//MaxServerID 最大的服务器ID号
	MaxServerID = 1<<27 - 1
)
