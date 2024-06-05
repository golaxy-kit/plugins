package variant

import (
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/util/binaryutil"
)

// Map map
type Map generic.SliceMap[Variant, Variant]

// Read implements io.Reader
func (v Map) Read(p []byte) (int, error) {
	bs := binaryutil.NewBigEndianStream(p)

	if err := bs.WriteUvarint(uint64(len(v))); err != nil {
		return bs.BytesWritten(), err
	}

	for i := range v {
		kv := &v[i]

		if _, err := binaryutil.ReadFrom(&bs, kv.K); err != nil {
			return bs.BytesWritten(), err
		}

		if _, err := binaryutil.ReadFrom(&bs, kv.V); err != nil {
			return bs.BytesWritten(), err
		}
	}

	return bs.BytesWritten(), nil
}

// Write implements io.Writer
func (v *Map) Write(p []byte) (int, error) {
	bs := binaryutil.NewBigEndianStream(p)

	l, err := bs.ReadUvarint()
	if err != nil {
		return bs.BytesRead(), err
	}

	*v = make([]generic.KV[Variant, Variant], l)

	for i := uint64(0); i < l; i++ {
		kv := &(*v)[i]

		if _, err := bs.WriteTo(&kv.K); err != nil {
			return bs.BytesRead(), err
		}

		if _, err := bs.WriteTo(&kv.V); err != nil {
			return bs.BytesRead(), err
		}
	}

	return bs.BytesRead(), nil
}

// Size 大小
func (v Map) Size() int {
	n := binaryutil.SizeofUvarint(uint64(len(v)))
	for i := range v {
		kv := &v[i]
		n += kv.K.Size()
		n += kv.V.Size()
	}
	return n
}

// TypeId 类型
func (Map) TypeId() TypeId {
	return TypeId_Map
}

// Indirect 原始值
func (v Map) Indirect() any {
	return v
}

// CastSliceMap 转换为SliceMap
func (v Map) CastSliceMap() generic.SliceMap[Variant, Variant] {
	return (generic.SliceMap[Variant, Variant])(v)
}
