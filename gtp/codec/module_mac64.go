package codec

import (
	"errors"
	"hash"
	"kit.golaxy.org/plugins/gtp"
	"kit.golaxy.org/plugins/gtp/binaryutil"
)

// MAC64Module MAC64模块
type MAC64Module struct {
	Hash       hash.Hash64 // hash(64bit)函数
	PrivateKey []byte      // 秘钥
	gcList     [][]byte    // GC列表
}

// PatchMAC 补充MAC
func (m *MAC64Module) PatchMAC(msgId gtp.MsgId, flags gtp.Flags, msgBuf []byte) (dst []byte, err error) {
	if m.Hash == nil {
		return nil, errors.New("setting Hash is nil")
	}

	m.Hash.Reset()
	bs := [2]byte{msgId, byte(flags)}
	m.Hash.Write(bs[:])
	m.Hash.Write(msgBuf)
	m.Hash.Write(m.PrivateKey)

	msgMAC := gtp.MsgMAC64{
		Data: msgBuf,
		MAC:  m.Hash.Sum64(),
	}

	buf := BytesPool.Get(msgMAC.Size())
	defer func() {
		if err == nil {
			m.gcList = append(m.gcList, buf)
		} else {
			BytesPool.Put(buf)
		}
	}()

	_, err = msgMAC.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// VerifyMAC 验证MAC
func (m *MAC64Module) VerifyMAC(msgId gtp.MsgId, flags gtp.Flags, msgBuf []byte) (dst []byte, err error) {
	if m.Hash == nil {
		return nil, errors.New("setting Hash is nil")
	}

	msgMAC := gtp.MsgMAC64{}

	_, err = msgMAC.Write(msgBuf)
	if err != nil {
		return nil, err
	}

	m.Hash.Reset()
	bs := [2]byte{msgId, byte(flags)}
	m.Hash.Write(bs[:])
	m.Hash.Write(msgMAC.Data)
	m.Hash.Write(m.PrivateKey)

	if m.Hash.Sum64() != msgMAC.MAC {
		return nil, ErrIncorrectMAC
	}

	return msgMAC.Data, nil
}

// SizeofMAC MAC大小
func (m *MAC64Module) SizeofMAC(msgLen int) int {
	return binaryutil.SizeofVarint(int64(msgLen)) + binaryutil.SizeofUint64()
}

// GC GC
func (m *MAC64Module) GC() {
	for i := range m.gcList {
		BytesPool.Put(m.gcList[i])
	}
	m.gcList = m.gcList[:0]
}