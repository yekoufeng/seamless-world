package main

import (
	"sync"
	"zeus/common"

	log "github.com/cihub/seelog"
)

type Props struct {
	Name       string
	Hp         uint32
	MaxHp      uint32
	Attack     uint32
	Defence    uint32
	State      uint32
	propsMutex *sync.Mutex
}

func (p *Props) CloneProps(ret *Props) {
	p.propsMutex.Lock()
	defer p.propsMutex.Unlock()

	*ret = *p
	ret.propsMutex = nil
}

// ReadValueFromStream 从Stream中读取属性
func (p *Props) ReadValueFromStream(name string, bs *common.ByteStream) {
	st, err := bs.ReadStr()
	if err != nil {
		return
	}

	if (name == "HP" || name == "MaxHP" || name == "Attack" || name == "Defence" || name == "State") && st != "uint32" {
		log.Error("type not match uint32 and ", st)
		return
	}

	p.propsMutex.Lock()
	defer p.propsMutex.Unlock()

	switch st {
	case "bool":
	case "int8":
	case "int16":
	case "int32":
	case "int64":
	case "uint8":
	case "uint16":
	case "uint32":
		v, err := bs.ReadUInt32()
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug(name, " ", v)
		if name == "HP" {
			p.Hp = v
		} else if name == "MaxHP" {
			p.MaxHp = v
		} else if name == "Attack" {
			p.Attack = v
		} else if name == "Defence" {
			p.Defence = v
		} else if name == "State" {
			p.State = v
		}
	case "uint64":
	case "float32":
	case "float64":
	case "string":
		str, err := bs.ReadStr()
		if err != nil {
			log.Error(err)
			return
		}
		if name == "Name" {
			p.Name = str
		}
	default:
		// proto struct
		// data, err := bs.ReadBytes()
		// if err != nil {
		// 	return nil, err
		// }
		// value, err := CreateInst(name)
		// if err != nil {
		// 	return nil, err
		// }
		// return value.(msgdef.IMsg).Unmarshal(data), nil

	}
}
