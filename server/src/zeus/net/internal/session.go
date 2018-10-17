package internal

import (
	"net"
	"zeus/net/internal/internal"
	"zeus/net/internal/types"
)

func NewSession(conn net.Conn, encryEnabled bool, msgCreator types.IMsgCreator) types.ISession {
	return internal.NewSession(conn, encryEnabled, msgCreator)
}
