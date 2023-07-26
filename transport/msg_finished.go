package transport

import (
	"kit.golaxy.org/plugins/transport/binaryutil"
)

// Finished消息标志位
const (
	Flag_EncryptOK  Flag = 1 << (iota + Flag_Customize) // 加密成功，在服务端发起的Finished消息携带
	Flag_AuthOK                                         // 鉴权成功，在服务端发起的Finished消息携带
	Flag_ContinueOK                                     // 断线重连成功，在服务端发起的Finished消息携带
)

// MsgFinished 握手结束，表示认可对端，可以开始传输数据
type MsgFinished struct {
	SendSeq uint32 // 服务端请求序号
	RecvSeq uint32 // 服务端响应序号
}

// Read implements io.Reader
func (m *MsgFinished) Read(p []byte) (int, error) {
	bs := binaryutil.NewByteStream(p)
	if err := bs.WriteUint32(m.SendSeq); err != nil {
		return 0, err
	}
	if err := bs.WriteUint32(m.RecvSeq); err != nil {
		return 0, err
	}
	return bs.BytesWritten(), nil
}

// Write implements io.Writer
func (m *MsgFinished) Write(p []byte) (int, error) {
	bs := binaryutil.NewByteStream(p)
	sendSeq, err := bs.ReadUint32()
	if err != nil {
		return 0, err
	}
	recvSeq, err := bs.ReadUint32()
	if err != nil {
		return 0, err
	}
	m.SendSeq = sendSeq
	m.RecvSeq = recvSeq
	return bs.BytesRead(), nil
}

// Size 消息大小
func (m *MsgFinished) Size() int {
	return binaryutil.SizeofUint32() + binaryutil.SizeofUint32() + binaryutil.SizeofUint32()
}

// MsgId 消息Id
func (MsgFinished) MsgId() MsgId {
	return MsgId_Finished
}

// Clone 克隆消息对象
func (m MsgFinished) Clone() Msg {
	return &m
}
