package transport

import "kit.golaxy.org/plugins/transport/binaryutil"

// MsgAuth 鉴权
type MsgAuth struct {
	Token      []byte // 令牌
	Extensions []byte // 扩展内容
}

func (m *MsgAuth) Read(p []byte) (int, error) {
	bs := binaryutil.NewByteStream(p)
	token, err := bs.ReadBytes()
	if err != nil {
		return 0, err
	}
	extensions, err := bs.ReadBytes()
	if err != nil {
		return 0, err
	}
	m.Token = token
	m.Extensions = extensions
	return bs.BytesRead(), nil
}

func (m *MsgAuth) Write(p []byte) (int, error) {
	bs := binaryutil.NewByteStream(p)
	if err := bs.WriteBytes(m.Token); err != nil {
		return 0, err
	}
	if err := bs.WriteBytes(m.Extensions); err != nil {
		return 0, err
	}
	return bs.BytesWritten(), nil
}

func (m *MsgAuth) Size() int {
	return binaryutil.SizeofBytes(m.Token) + binaryutil.SizeofBytes(m.Extensions)
}
