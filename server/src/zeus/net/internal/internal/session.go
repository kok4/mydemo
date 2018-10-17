package internal

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"time"
	"zeus/net/internal/internal/msgenc"
	"zeus/net/internal/internal/msghdl"
	"zeus/net/internal/types"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
)

func NewSession(conn net.Conn, encryEnabled bool, msgCreator types.IMsgCreator) *session {
	assert.True(msgCreator != nil, "msgCreator is nil")
	sess := &session{
		conn: conn,

		msgCreator: msgCreator,
		msgHdlr:    msghdl.New(msgCreator),

		sendBufC: make(chan []byte, 1024),

		hbTimerInterval: time.Duration(viper.GetInt("Config.HBTimer")) * time.Second,

		_encryEnabled: encryEnabled,
	}

	sess.hbTimer = time.NewTimer(sess.hbTimerInterval)
	sess.ctx, sess.ctxCancel = context.WithCancel(context.Background())
	return sess
}

// TODO: thread-safe

// session 代表一个网络连接
type session struct {
	conn net.Conn

	ctx       context.Context
	ctxCancel context.CancelFunc

	msgCreator types.IMsgCreator
	msgHdlr    *msghdl.MessageHandler

	sendBufC chan []byte

	hbTimer         *time.Timer
	hbTimerInterval time.Duration

	_isClosed     atomic.Bool
	_isError      atomic.Bool
	_encryEnabled bool

	// 关闭事件处理器
	onClosed func()

	// 会话验证功能，仅服务器有用。可以设置会话需要验证，第1个消息必须为验证消息。
	needVerify  bool        // 是否需要验证
	verifyMsgID types.MsgID // 验证消息ID
	isVerified  atomic.Bool // 是否已通过验证。由消息处理器设置。
}

// Start 验证完成
func (sess *session) Start() {
	assert.True(sess.msgCreator != nil)

	go sess.recvLoop() // 协程中会调用消息处理器
	go sess.sendLoop()
	if viper.GetBool("Config.HeartBeat") {
		go sess.hbLoop()
	}
}

// Send 发送消息
func (sess *session) Send(msg types.IMsg) {
	if msg == nil {
		return
	}

	if sess.IsClosed() {
		log.Warnf("Send after sess close %s %s %s",
			sess.conn.RemoteAddr(), reflect.TypeOf(msg),
			fmt.Sprintf("%s", msg)) // log是异步的，所以 msg 必须复制下。
		return
	}

	// Todo: 队列太长则断开
	msgBuf, err := msgenc.EncodeMsg(msg)
	if err != nil {
		log.Error("encode message error in Send(): ", err)
		return
	}

	sess.sendBufC <- msgBuf
}

// Touch 记录心跳状态
func (sess *session) ResetHb() {
	sess.hbTimer.Reset(sess.hbTimerInterval)
}

func (sess *session) hbLoop() {
	for {
		select {
		case <-sess.ctx.Done():
			sess.hbTimer.Stop()
			return
		case <-sess.hbTimer.C:
			log.Error("sess heart tick expired ", sess.conn.RemoteAddr())

			sess._isError.Store(true)
			sess.hbTimer.Stop()
			sess.Close()
			return
		}
	}
}

func (sess *session) recvLoop() {
	var err error
	assert.True(err == nil, "init error is not nil")

	for {
		select {
		case <-sess.ctx.Done():
			return
		default:
		} // select

		if err = sess.readAndHandleOneMsg(); err != nil {
			break
		}
	} // for

	assert.True(err != nil, "must be error")
	sess._isError.Store(true)
	if sess.IsClosed() {
		return
	}

	// 底层检测到连接断开，可能是客户端主动Close或客户端断网超时
	log.Errorf("read and handle message error: %s, %s", err, sess.conn.RemoteAddr())
	sess.Close()
	return
}

// Todo
//			if msg.Name() == "HeartBeat" {
//				sess.Send(&msgdef.HeartBeatResponse{})
//			} else if msg.Name() == "Ping" {
//				sess.Send(msg)
//			} else if msg.Name() != "HeartBeatResponse" {
//				sess.msgHdlr.HandleMsg(msg)
//			}
//			sess.ResetHb()

func (sess *session) sendLoop() {

	for {
		select {
		case <-sess.ctx.Done():
			return
		case buf := <-sess.sendBufC:
			msgBuf, err := msgenc.CompressAndEncrypt(buf, true, sess._encryEnabled)
			if err != nil {
				log.Error("compress and encrypt message error: ", err)
				continue
			}

			_, err = sess.conn.Write(msgBuf)
			if err != nil {
				sess._isError.Store(true)
				if sess.IsClosed() {
					return
				}
				log.Error("send message error ", err)
				sess.Close()
				return
			}
		}
	}
}

// Close 关闭.
// 所有发送完成后才关闭。或2s后强制关闭。
func (sess *session) Close() {
	if sess.IsClosed() {
		return
	}

	sess._isClosed.Store(true)

	sess.hbTimer.Stop()

	if sess.onClosed != nil {
		sess.onClosed()
	}

	if !sess._isError.Load() {
		go func() {
			closeTicker := time.NewTicker(100 * time.Millisecond)
			defer closeTicker.Stop()
			closeTimer := time.NewTimer(2 * time.Second)
			defer closeTimer.Stop()
			for {
				select {
				case <-closeTimer.C:
					sess.ctxCancel()
					sess.conn.Close()
					return
				case <-closeTicker.C:
					if len(sess.sendBufC) == 0 {
						sess.ctxCancel()
						sess.conn.Close()
						return
					}
				}
			}
		}()
	} else {
		sess.ctxCancel()
		sess.conn.Close()
	}
}

// RemoteAddr 远程地址
func (sess *session) RemoteAddr() string {
	return sess.conn.RemoteAddr().String()
}

// IsClosed 返回sess是否已经关闭
func (sess *session) IsClosed() bool {
	return sess._isClosed.Load()
}

// RegMsgProcFunc 注册消息处理函数.
func (sess *session) RegMsgProcFunc(msgID types.MsgID, procFun func(types.IMsg)) {
	sess.msgHdlr.RegMsgProcFunc(msgID, procFun)
}

func (sess *session) readAndHandleOneMsg() error {
	msgID, rawMsgBuf, err := readARQMsg(sess.conn)
	if err != nil {
		return err
	}

	assert.True(rawMsgBuf != nil, "rawMsg is nil")
	if !sess.needVerify {
		sess.msgHdlr.HandleRawMsg(msgID, rawMsgBuf)
		return nil
	}

	// 会话需要验证，第1个消息为验证请求消息
	if msgID != sess.verifyMsgID {
		msg := sess.msgCreator.NewMsg(msgID)
		vrf := sess.msgCreator.NewMsg(sess.verifyMsgID)
		return fmt.Errorf("need verify message ID %d(%s), but got %d(%s)",
			sess.verifyMsgID, reflect.TypeOf(vrf), msgID, reflect.TypeOf(msg))
	}
	sess.msgHdlr.HandleRawMsg(msgID, rawMsgBuf)

	if sess.isVerified.Load() {
		sess.needVerify = false // 已通过验证，不再需要了，
		return nil
	}

	return fmt.Errorf("session verification failed")
}

// SetVerified 设置会话已通过验证.
// thread-safe.
func (sess *session) SetVerified() {
	sess.isVerified.Store(true)
}

// SetVerifyMsgID 设置会话验证消息ID.
// 非线程安全，在Session.Start()之前设置。
func (sess *session) SetVerifyMsgID(verifyMsgID types.MsgID) {
	sess.needVerify = true
	sess.verifyMsgID = verifyMsgID
}

func (sess *session) SetOnClosed(onClosed func()) {
	sess.onClosed = onClosed
}
