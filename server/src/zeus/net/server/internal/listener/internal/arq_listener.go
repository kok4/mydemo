package internal

import (
	"context"
	"net"

	log "github.com/cihub/seelog"
	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/net/netutil"
)

// ARQListener 网络服务器监听
type ARQListener struct {
	protocal    string
	listener    net.Listener
	kcpListener net.Listener
	ctx         context.Context
	ctxCancel   context.CancelFunc

	maxConns int
}

func NewARQListener(protocal string, addr string, maxConns int) (*ARQListener, error) {
	l := &ARQListener{
		protocal:    protocal,
		listener:    nil,
		kcpListener: nil,
		maxConns:    maxConns,
	}
	l.ctx, l.ctxCancel = context.WithCancel(context.Background())

	err := l.listen(addr)
	return l, err
}

func (a *ARQListener) listen(addr string) error {
	var err error

	switch a.protocal {
	case "tcp":
		a.listener, err = net.Listen("tcp", addr)
	case "kcp":
		a.listener, err = kcp.Listen(addr)
	case "tcp+kcp":
		a.listener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		a.kcpListener, err = kcp.Listen(addr)
	default:
		panic("WRONG PROTOCAL")
	}
	if err != nil {
		return err
	}

	if a.maxConns > 0 {
		a.listener = netutil.LimitListener(a.listener, a.maxConns)
		//todo ？ maxConns优化
		if a.kcpListener != nil {
			a.kcpListener = netutil.LimitListener(a.kcpListener, a.maxConns)
		}
	}

	return nil
}

func (a *ARQListener) Close() {
	a.ctxCancel()
	a.listener.Close()
	if a.kcpListener != nil {
		a.kcpListener.Close()
	}
}

func (a *ARQListener) Run(connHandler IConnHandler) {
	// 同时开启kcp
	if a.kcpListener != nil {
		go func() {
			for {
				select {
				case <-a.ctx.Done():
					return
				default:
					{
						conn, err := a.kcpListener.Accept()
						if err != nil {
							log.Error("accept connection error ", err)
							continue
						}
						go connHandler.HandleConn(conn)
					}
				}
			}
		}()
	}

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			{
				conn, err := a.listener.Accept()
				if err != nil {
					log.Error("accept connection error ", err)
					continue
				}
				go connHandler.HandleConn(conn)
			}
		}
	}
}
