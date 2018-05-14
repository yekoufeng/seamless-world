package serializer

import (
	"reflect"
	"zeus/common"
	"zeus/msgdef"
)

const (
	_           = iota
	typeUint8   = 1
	typeUint16  = 2
	typeUint32  = 3
	typeUint64  = 4
	typeInt8    = 5
	typeInt16   = 6
	typeInt32   = 7
	typeInt64   = 8
	typeFloat32 = 9
	typeFloat64 = 10
	typeString  = 11
	typeBytes   = 12
	typeBool    = 13
	typeProto   = 14
)

// Serialize 序列化
func Serialize(args ...interface{}) []byte {
	size := getSize(args...)
	data := make([]byte, size)
	bw := common.NewByteStream(data)

	for _, arg := range args {
		var err error

		switch arg.(type) {
		case uint8:
			err = bw.WriteByte(typeUint8)
			err = bw.WriteByte(arg.(uint8))
		case uint16:
			err = bw.WriteByte(typeUint16)
			err = bw.WriteUInt16(arg.(uint16))
		case uint32:
			err = bw.WriteByte(typeUint32)
			err = bw.WriteUInt32(arg.(uint32))
		case uint64:
			err = bw.WriteByte(typeUint64)
			err = bw.WriteUInt64(arg.(uint64))
		case int8:
			err = bw.WriteByte(typeInt8)
			err = bw.WriteInt8(arg.(int8))
		case int16:
			err = bw.WriteByte(typeInt16)
			err = bw.WriteInt16(arg.(int16))
		case int32:
			err = bw.WriteByte(typeInt32)
			err = bw.WriteInt32(arg.(int32))
		case int64:
			err = bw.WriteByte(typeInt64)
			err = bw.WriteInt64(arg.(int64))
		case float32:
			err = bw.WriteByte(typeFloat32)
			err = bw.WriteFloat32(arg.(float32))
		case float64:
			err = bw.WriteByte(typeFloat64)
			err = bw.WriteFloat64(arg.(float64))
		case string:
			err = bw.WriteByte(typeString)
			err = bw.WriteStr(arg.(string))
		case []byte:
			err = bw.WriteByte(typeBytes)
			err = bw.WriteBytes(arg.([]byte))
		case bool:
			err = bw.WriteByte(typeBool)
			err = bw.WriteBool(arg.(bool))
		default:
			// proto消息序列化格式:
			// 数据类型 | 消息号 | 消息内容
			//    1         2     2(字节流长度) + 字节流
			iM := arg.(msgdef.IMsg)
			err = bw.WriteByte(typeProto)
			if err != nil {
				panic(err)
			}

			var msgID uint16
			msgID, err = msgdef.GetMsgDef().GetMsgIDByName(reflect.TypeOf(arg).Elem().Name())
			err = bw.WriteUInt16(msgID)
			if err != nil {
				panic(err)
			}

			iMData := make([]byte, iM.Size())
			if _, err = iM.MarshalTo(iMData); err != nil {
				panic(err)
			}
			err = bw.WriteBytes(iMData)
		}

		if err != nil {
			panic(err)
		}
	}

	return data
}

// UnSerialize 反序列化
func UnSerialize(data []byte) []interface{} {
	ret := make([]interface{}, 0, 1)
	br := common.NewByteStream(data)

	for br.ReadEnd() {
		var err error
		var typ uint8
		var v interface{}
		if typ, err = br.ReadByte(); err != nil {
			panic(err)
		}

		switch typ {
		case typeUint8:
			v, err = br.ReadByte()
		case typeUint16:
			v, err = br.ReadUInt16()
		case typeUint32:
			v, err = br.ReadUInt32()
		case typeUint64:
			v, err = br.ReadUInt64()
		case typeInt8:
			v, err = br.ReadInt8()
		case typeInt16:
			v, err = br.ReadInt16()
		case typeInt32:
			v, err = br.ReadInt32()
		case typeInt64:
			v, err = br.ReadInt64()
		case typeFloat32:
			v, err = br.ReadFloat32()
		case typeFloat64:
			v, err = br.ReadFloat64()
		case typeString:
			v, err = br.ReadStr()
		case typeBytes:
			v, err = br.ReadBytes()
		case typeBool:
			v, err = br.ReadBool()
		case typeProto:
			msgID, err := br.ReadUInt16()
			if err != nil {
				panic(err)
			}
			_, msgContent, err := msgdef.GetMsgDef().GetMsgInfo(msgID)
			if err != nil {
				panic(err)
			}
			if iMData, err := br.ReadBytes(); err == nil {
				err = msgContent.Unmarshal(iMData)
				v = msgContent
			} else {
				panic(err)
			}
		default:
		}

		if err == nil {
			ret = append(ret, v)
		} else {
			panic(err)
		}
	}

	if len(ret) == 0 {
		return nil
	}

	return ret
}

// 每个参数需要固定1个字节表示数据类型
// 每种类型的长度不固定
func getSize(args ...interface{}) int {
	size := 0
	for _, arg := range args {
		//类型1个字节
		size++

		switch arg.(type) {
		case uint8, int8, bool:
			size++
		case uint16, int16:
			size += 2
		case uint32, int32, float32:
			size += 4
		case uint64, int64, float64:
			size += 8
		case string:
			// 字符串需要2个字节标识长度+本身的长度
			size += 2
			size += len(arg.(string))
		case []byte:
			size += 2
			size += len(arg.([]byte))
		default:
			iM := arg.(msgdef.IMsg)
			// proto需要2个字节的消息号+2个字节标识消息长度+消息本身的长度
			size += 4
			size += iM.Size()
		}
	}

	return size
}
