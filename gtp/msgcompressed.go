package gtp

import (
	"kit.golaxy.org/plugins/util/binaryutil"
)

// MsgCompressed 压缩消息
type MsgCompressed struct {
	Data         []byte
	OriginalSize int64
}

// Read implements io.Reader
func (m MsgCompressed) Read(p []byte) (int, error) {
	bs := binaryutil.NewBigEndianStream(p)
	if err := bs.WriteBytes(m.Data); err != nil {
		return bs.BytesWritten(), err
	}
	if err := bs.WriteVarint(m.OriginalSize); err != nil {
		return bs.BytesWritten(), err
	}
	return bs.BytesWritten(), nil
}

// Write implements io.Writer
func (m *MsgCompressed) Write(p []byte) (int, error) {
	bs := binaryutil.NewBigEndianStream(p)
	data, err := bs.ReadBytesRef()
	if err != nil {
		return bs.BytesRead(), err
	}
	originalSize, err := bs.ReadVarint()
	if err != nil {
		return bs.BytesRead(), err
	}
	m.Data = data
	m.OriginalSize = originalSize
	return bs.BytesRead(), nil
}

// Size 大小
func (m MsgCompressed) Size() int {
	return binaryutil.SizeofBytes(m.Data) + binaryutil.SizeofVarint(m.OriginalSize)
}
