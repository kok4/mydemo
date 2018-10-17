package internal

import "net"

type IConnHandler interface {
	HandleConn(net.Conn)
}
