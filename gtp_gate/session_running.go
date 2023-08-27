package gtp_gate

import (
	"bytes"
	"errors"
	"fmt"
	"kit.golaxy.org/golaxy/util"
	"kit.golaxy.org/plugins/gtp"
	"kit.golaxy.org/plugins/gtp/codec"
	"kit.golaxy.org/plugins/gtp/transport"
	"kit.golaxy.org/plugins/internal"
	"kit.golaxy.org/plugins/logger"
	"math/rand"
	"net"
	"time"
)

// Init 初始化
func (s *_GtpSession) Init(conn net.Conn, encoder codec.IEncoder, decoder codec.IDecoder, token string) (sendSeq, recvSeq uint32) {
	s.Lock()
	defer s.Unlock()

	// 初始化消息收发器
	s.transceiver.Conn = conn
	s.transceiver.Encoder = encoder
	s.transceiver.Decoder = decoder
	s.transceiver.Timeout = s.gate.options.IOTimeout

	buff := &transport.SequencedBuffer{}
	buff.Reset(rand.Uint32(), rand.Uint32(), s.gate.options.IOBufferCap)

	s.transceiver.Buffer = buff

	// 初始化刷新通知channel
	s.renewChan = make(chan struct{}, 1)

	// 初始化token
	s.token = token

	return buff.SendSeq(), buff.RecvSeq()
}

// Renew 刷新
func (s *_GtpSession) Renew(conn net.Conn, remoteRecvSeq uint32) (sendSeq, recvSeq uint32, err error) {
	s.Lock()

	defer func() {
		s.Unlock()

		select {
		case s.renewChan <- struct{}{}:
		default:
		}
	}()

	// 刷新链路
	return s.transceiver.Renew(conn, remoteRecvSeq)
}

// PauseIO 暂停收发消息
func (s *_GtpSession) PauseIO() {
	s.transceiver.Pause()
}

// ContinueIO 继续收发消息
func (s *_GtpSession) ContinueIO() {
	s.transceiver.Continue()
}

// Run 运行（会话的主线程）
func (s *_GtpSession) Run() {
	defer func() {
		if panicErr := util.Panic2Err(recover()); panicErr != nil {
			logger.Errorf(s.gate.ctx, "session %q panicked, %s", s.GetId(), fmt.Errorf("panicked: %w", panicErr))
		}

		// 调整会话状态为已过期
		s.SetState(SessionState_Death)

		// 关闭连接和清理数据
		if s.transceiver.Conn != nil {
			s.transceiver.Conn.Close()
		}
		s.transceiver.Clean()

		// 删除会话
		s.gate.DeleteSession(s.GetId())

		logger.Debugf(s.gate.ctx, "session %q shutdown, conn %q -> %q", s.GetId(), s.GetLocalAddr(), s.GetRemoteAddr())
	}()

	logger.Debugf(s.gate.ctx, "session %q started, conn %q -> %q", s.GetId(), s.GetLocalAddr(), s.GetRemoteAddr())

	// 启动发送数据的线程
	if s.options.SendDataChan != nil {
		go func() {
			for {
				select {
				case data := <-s.options.SendDataChan:
					if err := s.SendData(data); err != nil {
						logger.Errorf(s.gate.ctx, "session %q fetch data from the send data channel for sending failed, %s", s.GetId(), err)
					}
				case <-s.Done():
					return
				}
			}
		}()
	}

	// 启动发送自定义事件的线程
	if s.options.SendEventChan != nil {
		go func() {
			for {
				select {
				case event := <-s.options.SendEventChan:
					if err := s.SendEvent(event); err != nil {
						logger.Errorf(s.gate.ctx, "session %q fetch event from the send event channel for sending failed, %s", s.GetId(), err)
					}
				case <-s.Done():
					return
				}
			}
		}()
	}

	pinged := false
	var timeout time.Time

	// 调整会话状态为活跃
	s.SetState(SessionState_Active)

	for {
		// 非活跃状态，检测超时时间
		if s.state == SessionState_Inactive {
			if time.Now().After(timeout) {
				s.cancel()
			}
		}

		// 检测会话是否已关闭
		select {
		case <-s.Done():
			return
		default:
		}

		// 分发消息事件
		if err := s.dispatcher.Dispatching(); err != nil {
			logger.Debugf(s.gate.ctx, "session %q dispatching event failed, %s", s.GetId(), err)

			// 网络io超时，触发心跳检测，向对方发送ping
			if errors.Is(err, transport.ErrTimeout) {
				if !pinged {
					logger.Debugf(s.gate.ctx, "session %q send ping", s.GetId())

					s.ctrl.SendPing()
					pinged = true
				} else {
					logger.Debugf(s.gate.ctx, "session %q no pong received", s.GetId())

					// 未收到对方回复pong或其他消息事件，再次网络io超时，调整会话状态不活跃
					if s.SetState(SessionState_Inactive) {
						timeout = time.Now().Add(s.gate.options.SessionInactiveTimeout)
					}
				}
				continue
			}

			// 其他网络io类错误，调整会话状态不活跃
			if errors.Is(err, transport.ErrNetIO) {
				if s.SetState(SessionState_Inactive) {
					timeout = time.Now().Add(s.gate.options.SessionInactiveTimeout)
				}

				func() {
					timer := time.NewTimer(10 * time.Second)
					defer timer.Stop()

					select {
					case <-timer.C:
						return
					case <-s.renewChan:
						// 发送缓存的消息
						transport.Retry{
							Transceiver: &s.transceiver,
							Times:       s.gate.options.IORetryTimes,
						}.Send(s.transceiver.Resend())
						// 重置ping状态
						pinged = false
						return
					case <-s.Done():
						return
					}
				}()

				logger.Debugf(s.gate.ctx, "session %q retry dispatching event, conn %q -> %q", s.GetId(), s.GetLocalAddr(), s.GetRemoteAddr())
				continue
			}

			continue
		}

		// 没有错误，或非网络io类错误，重置ping状态
		pinged = false
		// 调整会话状态活跃
		s.SetState(SessionState_Active)
	}
}

// SetState 调整会话状态
func (s *_GtpSession) SetState(state SessionState) bool {
	old := s.state

	if old == state {
		return false
	}

	s.Lock()
	s.state = state
	s.Unlock()

	logger.Debugf(s.gate.ctx, "session %q state %q => %q", s.GetId(), old, state)

	for i := range s.gate.options.SessionStateChangedHandlers {
		handler := s.gate.options.SessionStateChangedHandlers[i]
		if handler == nil {
			continue
		}
		err := internal.CallVoid(func() { handler(s, old, state) })
		if err != nil {
			logger.Errorf(s.gate.ctx, "session %q state changed handler error: %s", s.GetId(), err)
		}
	}

	for i := range s.options.StateChangedHandlers {
		handler := s.options.StateChangedHandlers[i]
		if handler == nil {
			continue
		}
		err := internal.CallVoid(func() { handler(old, state) })
		if err != nil {
			logger.Errorf(s.gate.ctx, "session %q state changed handler error: %s", s.GetId(), err)
		}
	}

	return true
}

// EventHandler 接收自定义事件的处理器
func (s *_GtpSession) EventHandler(event transport.Event[gtp.Msg]) error {
	if s.options.RecvEventChan != nil {
		select {
		case s.options.RecvEventChan <- event.Clone():
		default:
			logger.Errorf(s.gate.ctx, "session %q receive event channel is full", s.GetId())
		}
	}

	for i := range s.gate.options.SessionRecvEventHandlers {
		handler := s.gate.options.SessionRecvEventHandlers[i]
		if handler == nil {
			continue
		}
		err := internal.Call(func() error { return handler(s, event) })
		if err == nil || !errors.Is(err, transport.ErrUnexpectedMsg) {
			return err
		}
	}

	for i := range s.options.RecvEventHandlers {
		handler := s.options.RecvEventHandlers[i]
		if handler == nil {
			continue
		}
		err := internal.Call(func() error { return handler(event) })
		if err == nil || !errors.Is(err, transport.ErrUnexpectedMsg) {
			return err
		}
	}

	return transport.ErrUnexpectedMsg
}

// PayloadHandler Payload消息事件处理器
func (s *_GtpSession) PayloadHandler(event transport.Event[*gtp.MsgPayload]) error {
	if s.options.RecvDataChan != nil {
		select {
		case s.options.RecvDataChan <- bytes.Clone(event.Msg.Data):
		default:
			logger.Errorf(s.gate.ctx, "session %q receive data channel is full", s.GetId())
		}
	}

	for i := range s.gate.options.SessionRecvDataHandlers {
		handler := s.gate.options.SessionRecvDataHandlers[i]
		if handler == nil {
			continue
		}
		err := internal.Call(func() error { return handler(s, event.Msg.Data) })
		if err == nil || !errors.Is(err, transport.ErrUnexpectedMsg) {
			return err
		}
	}

	for i := range s.options.RecvDataHandlers {
		handler := s.options.RecvDataHandlers[i]
		if handler == nil {
			continue
		}
		err := internal.Call(func() error { return handler(event.Msg.Data) })
		if err == nil || !errors.Is(err, transport.ErrUnexpectedMsg) {
			return err
		}
	}

	return transport.ErrUnexpectedMsg
}

// HeartbeatHandler Heartbeat消息事件处理器
func (s *_GtpSession) HeartbeatHandler(event transport.Event[*gtp.MsgHeartbeat]) error {
	if event.Flags.Is(gtp.Flag_Ping) {
		logger.Debugf(s.gate.ctx, "session %q receive ping", s.GetId())
	} else {
		logger.Debugf(s.gate.ctx, "session %q receive pong", s.GetId())
	}
	return nil
}
