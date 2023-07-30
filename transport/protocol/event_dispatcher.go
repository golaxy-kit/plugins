package protocol

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"kit.golaxy.org/plugins/internal"
	"kit.golaxy.org/plugins/transport"
)

var (
	ErrHandlerNotRegistered = errors.New("handler not registered") // 消息处理器未注册
	ErrUnexpectedMsg        = errors.New("unexpected msg")         // 收到非预期的消息
)

type (
	EventHandler = func(Event[transport.Msg]) error     // 消息事件处理器
	ErrorHandler = func(ctx context.Context, err error) // 错误处理器
)

// EventDispatcher 消息事件分发器
type EventDispatcher struct {
	Transceiver   *Transceiver   // 消息事件收发器
	RetryTimes    int            // 网络io超时时的重试次数
	EventHandlers []EventHandler // 消息事件处理器
	ErrorHandler  ErrorHandler   // 错误处理器
}

// Run 运行
func (d *EventDispatcher) Run(ctx context.Context) {
	if d.Transceiver == nil {
		return
	}

	defer d.Transceiver.GC()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		d.Transceiver.GC()

		e, err := d.retryRecv(d.Transceiver.Recv())
		if err != nil {
			if d.ErrorHandler != nil {
				internal.CallVoid(func() { d.ErrorHandler(ctx, err) })
			}
			continue
		}

		handled := false

		for i := range d.EventHandlers {
			if err = internal.Call(func() error { return d.EventHandlers[i](e) }); err != nil {
				if errors.Is(err, ErrUnexpectedMsg) {
					continue
				}
				if d.ErrorHandler != nil {
					internal.CallVoid(func() { d.ErrorHandler(ctx, err) })
				}
			}
			handled = true
			break
		}

		if !handled {
			if d.ErrorHandler != nil {
				internal.CallVoid(func() { d.ErrorHandler(ctx, fmt.Errorf("%w: %d", ErrHandlerNotRegistered, e.Msg.MsgId())) })
			}
		}
	}
}

func (d *EventDispatcher) retryRecv(e Event[transport.Msg], err error) (Event[transport.Msg], error) {
	return Retry{
		Transceiver: d.Transceiver,
		Times:       d.RetryTimes,
	}.Recv(e, err)
}
