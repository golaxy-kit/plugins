package tcp

import (
	"errors"
	"fmt"
	"kit.golaxy.org/plugins/transport"
	"kit.golaxy.org/plugins/transport/codec"
	"net"
	"os"
)

// MsgBlock 消息块
type MsgBlock[T transport.Msg] struct {
	Flags transport.Flags // 标志位
	Msg   T               // 消息
}

// RstError Rst错误提示
type RstError struct {
	Code    transport.Code // 错误码
	Message string         // 错误信息
}

// Error 错误信息
func (e *RstError) Error() string {
	return fmt.Sprintf("(%d)%s", e.Code, e.Message)
}

// Transceiver 消息收发器
type Transceiver struct {
	Conn       net.Conn       // 网络连接
	Encoder    codec.IEncoder // 消息包编码器
	Decoder    codec.IDecoder // 消息包解码器
	RetryTimes int            // io超时重试次数
}

// Send 发送消息
func (t Transceiver) Send(mb MsgBlock[transport.Msg]) error {
	if t.Conn == nil {
		return errors.New("conn is nil")
	}

	if t.Encoder == nil {
		return errors.New("encoder is nil")
	}

	if err := t.Encoder.Stuff(mb.Flags, mb.Msg); err != nil {
		return fmt.Errorf("stuff send mb failed, %s", err)
	}

	var retries int
retry:
	if _, err := t.Encoder.WriteTo(t.Conn); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			if retries < t.RetryTimes {
				retries++
				goto retry
			}
		}
		return fmt.Errorf("send mb-packet failed, %s", err)
	}

	return nil
}

// SendRst 发送Rst
func (t Transceiver) SendRst(err error) {
	if t.Conn == nil || t.Encoder == nil {
		return
	}

	msg := &transport.MsgRst{}

	rstErr, ok := err.(*RstError)
	if ok {
		msg.Code = rstErr.Code
		msg.Message = rstErr.Message
	} else {
		msg.Code = transport.Code_Reject
		msg.Message = rstErr.Message
	}

	if err := t.Encoder.Stuff(transport.Flags_None, msg); err != nil {
		return
	}

	if _, err := t.Encoder.WriteTo(t.Conn); err != nil {
		return
	}
}

// Recv 接收单个消息
func (t Transceiver) Recv() (MsgBlock[transport.Msg], error) {
	if t.Conn == nil {
		return MsgBlock[transport.Msg]{}, errors.New("conn is nil")
	}

	if t.Decoder == nil {
		return MsgBlock[transport.Msg]{}, errors.New("decoder is nil")
	}

	for {
		var recvMP transport.MsgPacket

		if err := t.Decoder.Fetch(func(mp transport.MsgPacket) { recvMP = mp }); err != nil {
			if !errors.Is(err, codec.ErrEmptyCache) {
				return MsgBlock[transport.Msg]{}, fmt.Errorf("fetch recv msg-packet failed, %s", err)
			}
		} else {
			return MsgBlock[transport.Msg]{
				Flags: recvMP.Head.Flags,
				Msg:   recvMP.Msg,
			}, nil
		}

		var retries int
	retry:
		if _, err := t.Decoder.ReadFrom(t.Conn); err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if retries < t.RetryTimes {
					retries++
					goto retry
				}
			}
			return MsgBlock[transport.Msg]{}, fmt.Errorf("recv msg-packet failed, %s", err)
		}
	}
}

// MultiRecv 接收多个消息
func (t Transceiver) MultiRecv(fun func(MsgBlock[transport.Msg]) bool) error {
	if fun == nil {
		return errors.New("fun is nil")
	}

	if t.Conn == nil {
		return errors.New("conn is nil")
	}

	if t.Decoder == nil {
		return errors.New("decoder is nil")
	}

	for {
		var b bool

		err := t.Decoder.MultiFetch(func(mp transport.MsgPacket) bool {
			b = fun(MsgBlock[transport.Msg]{
				Flags: mp.Head.Flags,
				Msg:   mp.Msg,
			})
			return b
		})
		if err != nil {
			if !errors.Is(err, codec.ErrEmptyCache) {
				return fmt.Errorf("fetch recv msg-packet failed, %s", err)
			}
		} else {
			if !b {
				return nil
			}
		}

		var retries int
	retry:
		if _, err := t.Decoder.ReadFrom(t.Conn); err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if retries < t.RetryTimes {
					retries++
					goto retry
				}
			}
			return fmt.Errorf("recv msg-packet failed, %s", err)
		}
	}
}
