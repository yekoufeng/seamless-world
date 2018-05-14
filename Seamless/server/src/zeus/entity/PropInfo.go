package entity

import (
	"fmt"
	"zeus/common"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// PropInfo 属性相关
type PropInfo struct {
	value interface{}

	syncFlag bool
	dbFlag   bool

	def *PropDef
}

func newPropInfo(def *PropDef) *PropInfo {
	prop := &PropInfo{
		value:    nil,
		syncFlag: false,
		dbFlag:   false,
		def:      def,
	}
	prop.init()
	return prop
}

func (p *PropInfo) init() {
	if p.def == nil {
		log.Error("属性初始化失败, Def为空")
		return
	}

	switch p.def.TypeName {
	case "bool":
		p.value = false
	case "int8":
		p.value = int8(0)
	case "int16":
		p.value = int16(0)
	case "int32":
		p.value = int32(0)
	case "int64":
		p.value = int64(0)
	case "uint8":
		p.value = uint8(0)
	case "uint16":
		p.value = uint16(0)
	case "uint32":
		p.value = uint32(0)
	case "uint64":
		p.value = uint64(0)
	case "float32":
		p.value = float32(0)
	case "float64":
		p.value = float64(0)
	case "string":
		p.value = ""
	default:
		var err error
		p.value, err = p.def.CreateInst()
		if err != nil {
			log.Error(err, p.def)
			return
		}
	}
}

// GetValueStreamSize 获取某个属性需要的尺寸
func (p *PropInfo) GetValueStreamSize() int {
	s := 0
	st := p.def.TypeName

	switch st {
	case "int8", "uint8", "bool":
		s = 1
	case "int16", "uint16":
		s = 2
	case "int32", "uint32":
		s = 4
	case "int64", "uint64":
		s = 8
	case "float32":
		s = 4
	case "float64":
		s = 8
	case "string":
		s = len(p.value.(string)) + 2
	default:
		// proto struct
		// 本身序列化的长度+写入字节流时需要额外2个字节标识自身的长度
		if iM, ok := p.value.(msgdef.IMsg); ok {
			s = iM.Size() + 2
		} else {
			log.Error("Convert proto struct failed", st)
		}
	}

	return s + len(st) + 2
}

// WriteValueToStream 把属性值加入到ByteStream中
func (p *PropInfo) WriteValueToStream(bs *common.ByteStream) error {
	err := bs.WriteStr(p.def.TypeName)
	if err != nil {
		return err
	}

	st := p.def.TypeName

	switch st {
	case "bool":
		err = bs.WriteBool(bool(p.value.(bool)))
	case "int8":
		err = bs.WriteByte(byte(p.value.(int8)))
	case "int16":
		err = bs.WriteUInt16(uint16(p.value.(int16)))
	case "int32":
		err = bs.WriteUInt32(uint32(p.value.(int32)))
	case "int64":
		err = bs.WriteUInt64(uint64(p.value.(int64)))
	case "uint8":
		err = bs.WriteByte(p.value.(byte))
	case "uint16":
		err = bs.WriteUInt16(p.value.(uint16))
	case "uint32":
		err = bs.WriteUInt32(p.value.(uint32))
	case "uint64":
		err = bs.WriteUInt64(p.value.(uint64))
	case "string":
		err = bs.WriteStr(p.value.(string))
	case "float32":
		err = bs.WriteFloat32(p.value.(float32))
	case "float64":
		err = bs.WriteFloat64(p.value.(float64))
	default:
		// proto struct
		if iM, ok := p.value.(msgdef.IMsg); ok {
			data := make([]byte, iM.Size())
			if _, err = iM.MarshalTo(data); err != nil {
				return err
			}
			err = bs.WriteBytes(data)
		} else {
			log.Error("Convert proto struct failed", st)
		}
	}

	return err
}

// ReadValueFromStream 从Stream中读取属性
func (p *PropInfo) ReadValueFromStream(bs *common.ByteStream) error {
	st, err := bs.ReadStr()
	if err != nil {
		return err
	}

	switch st {
	case "bool":
		v, err := bs.ReadBool()
		if err != nil {
			return err
		}
		p.value = bool(v)
	case "int8":
		v, err := bs.ReadByte()
		if err != nil {
			return err
		}
		p.value = int8(v)
	case "int16":
		v, err := bs.ReadUInt16()
		if err != nil {
			return err
		}
		p.value = int16(v)
	case "int32":
		v, err := bs.ReadUInt32()
		if err != nil {
			return err
		}
		p.value = int32(v)
	case "int64":
		v, err := bs.ReadUInt64()
		if err != nil {
			return err
		}
		p.value = int64(v)
	case "uint8":
		v, err := bs.ReadByte()
		if err != nil {
			return err
		}
		p.value = v
	case "uint16":
		v, err := bs.ReadUInt16()
		if err != nil {
			return err
		}
		p.value = v
	case "uint32":
		v, err := bs.ReadUInt32()
		if err != nil {
			return err
		}
		p.value = v
	case "uint64":
		v, err := bs.ReadUInt64()
		if err != nil {
			return err
		}
		p.value = v
	case "float32":
		v, err := bs.ReadFloat32()
		if err != nil {
			return err
		}
		p.value = v
	case "float64":
		v, err := bs.ReadFloat64()
		if err != nil {
			return err
		}
		p.value = v
	case "string":
		v, err := bs.ReadStr()
		if err != nil {
			return err
		}
		p.value = v
	default:
		// proto struct
		data, err := bs.ReadBytes()
		if err != nil {
			return err
		}
		p.value, err = p.def.CreateInst()
		if err != nil {
			return err
		}
		return p.value.(msgdef.IMsg).Unmarshal(data)

	}

	return nil
}

// PackValue 打包value 给Redis用
func (p *PropInfo) PackValue() []byte {
	switch p.def.TypeName {
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool", "string":
		return []byte(fmt.Sprintf("%v", p.value))
	default:
		iM, ok := p.value.(msgdef.IMsg)
		if ok {
			data := make([]byte, iM.Size())
			iM.MarshalTo(data)
			return data
		}

		return nil
	}
}

// UnPackValue 从Redis中恢复Value
func (p *PropInfo) UnPackValue(data interface{}) {
	switch p.def.TypeName {
	case "bool":
		v, _ := redis.Bool(data, nil)
		p.value = bool(v)
	case "int8":
		v, _ := redis.Int(data, nil)
		p.value = int8(v)
	case "int16":
		v, _ := redis.Int(data, nil)
		p.value = int16(v)
	case "int32":
		v, _ := redis.Int(data, nil)
		p.value = int32(v)
	case "int64":
		p.value, _ = redis.Int64(data, nil)
	case "uint8":
		v, _ := redis.Uint64(data, nil)
		p.value = uint8(v)
	case "uint16":
		v, _ := redis.Uint64(data, nil)
		p.value = uint16(v)
	case "uint32":
		v, _ := redis.Uint64(data, nil)
		p.value = uint32(v)
	case "uint64":
		p.value, _ = redis.Uint64(data, nil)
	case "float32":
		v, _ := redis.Float64(data, nil)
		p.value = float32(v)
	case "float64":
		p.value, _ = redis.Float64(data, nil)
	case "string":
		p.value, _ = redis.String(data, nil)
	default:
		if data != nil {
			iM, ok := p.value.(msgdef.IMsg)
			if ok {
				if err := iM.Unmarshal(data.([]byte)); err != nil {
					log.Error(err)
				}
			} else {
				log.Warn("Unsupport prop type", p.def.TypeName)
			}
		}
	}
}
