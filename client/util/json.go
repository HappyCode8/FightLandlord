package util

import (
	"github.com/json-iterator/go"
	"unsafe"
)

func init() {
	// 编码器，将[]uint8类型的数据（即[]byte）序列化为JSON字符串
	jsoniter.RegisterTypeEncoderFunc("[]uint8", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		t := *((*[]byte)(ptr))
		stream.WriteString(string(t))
	}, nil)
	// 解码器，将JSON字符串反序列化为[]uint8类型的数据
	jsoniter.RegisterTypeDecoderFunc("[]uint8", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		str := iter.ReadString()
		*((*[]byte)(ptr)) = []byte(str)
	})
}

func Marshal(v interface{}) []byte {
	data, _ := jsoniter.Marshal(v)
	return data
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
