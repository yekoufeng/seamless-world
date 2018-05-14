package entitydef

import "zeus/iserver"

// ItemDef 自动生成的属性包装代码
type ItemDef struct {
	ip iserver.IEntityProps
}

// SetPropsSetter 设置接口
func (p *ItemDef) SetPropsSetter(ip iserver.IEntityProps) {
	p.ip = ip
}

// Setbaseid 设置 baseid
func (p *ItemDef) Setbaseid(v uint32) {
	p.ip.SetProp("baseid", v)
}

// SetbaseidDirty 设置baseid被修改
func (p *ItemDef) SetbaseidDirty() {
	p.ip.PropDirty("baseid")
}

// Getbaseid 获取 baseid
func (p *ItemDef) Getbaseid() uint32 {
	v := p.ip.GetProp("baseid")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// Setnum 设置 num
func (p *ItemDef) Setnum(v uint32) {
	p.ip.SetProp("num", v)
}

// SetnumDirty 设置num被修改
func (p *ItemDef) SetnumDirty() {
	p.ip.PropDirty("num")
}

// Getnum 获取 num
func (p *ItemDef) Getnum() uint32 {
	v := p.ip.GetProp("num")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

// Setreducedam 设置 reducedam
func (p *ItemDef) Setreducedam(v uint32) {
	p.ip.SetProp("reducedam", v)
}

// SetreducedamDirty 设置reducedam被修改
func (p *ItemDef) SetreducedamDirty() {
	p.ip.PropDirty("reducedam")
}

// Getreducedam 获取 reducedam
func (p *ItemDef) Getreducedam() uint32 {
	v := p.ip.GetProp("reducedam")
	if v == nil {
		return uint32(0)
	}

	return v.(uint32)
}

type IItemDef interface {
	Setbaseid(v uint32)
	SetbaseidDirty()
	Getbaseid() uint32
	Setnum(v uint32)
	SetnumDirty()
	Getnum() uint32
	Setreducedam(v uint32)
	SetreducedamDirty()
	Getreducedam() uint32
}
