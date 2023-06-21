package protocol

import (
	"errors"
	"golang.org/x/net/context"
	"kit.golaxy.org/plugins/transport"
	"kit.golaxy.org/plugins/transport/codec"
	"net"
	"time"
)

var (
	ErrHandlerNotRegistered = errors.New("handler not registered") // 消息处理句柄未注册
	ErrRecvUnexpectedMsg    = errors.New("recv unexpected msg")    // 收到非预期的消息
)

// ErrorHandler 错误句柄
type ErrorHandler = func(err error) bool

// Handler 消息句柄
type Handler interface {
	Recv(Event[transport.Msg]) error
}

// Dispatcher 消息事件分发器
type Dispatcher struct {
	Conn     net.Conn                    // 网络连接
	Decoder  codec.IDecoder              // 消息包解码器
	Timeout  time.Duration               // io超时时间
	Handlers map[transport.MsgId]Handler // 消息处理句柄
}

// Run 运行
func (d *Dispatcher) Run(ctx context.Context, errorHandler ErrorHandler) {
	trans := Transceiver{
		Conn:    d.Conn,
		Decoder: d.Decoder,
		Timeout: d.Timeout,
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		e, err := trans.Recv()
		if err != nil {
			if errorHandler != nil {
				if !errorHandler(err) {
					return
				}
			}
			continue
		}

		var handler Handler

		if d.Handlers != nil {
			handler = d.Handlers[e.Msg.MsgId()]
		}

		if handler == nil {
			if errorHandler != nil {
				if !errorHandler(ErrHandlerNotRegistered) {
					return
				}
			}
			continue
		}

		err = handler.Recv(e)
		if err != nil {
			if errorHandler != nil {
				if !errorHandler(err) {
					return
				}
			}
			continue
		}
	}
}
