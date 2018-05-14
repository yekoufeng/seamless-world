package main

import (
	"fmt"
	"sort"
	"strings"
)

// Template 代码生成模版
type Template struct {
	name         string
	typeInfo     map[string]string
	needAddProto bool
	interfaceStr string
}

// NewTemplate 生成新的模版工具
func NewTemplate(name string) *Template {
	t := &Template{}
	t.name = name
	t.typeInfo = make(map[string]string)
	t.needAddProto = false
	return t
}

func (t *Template) String() string {
	str := t.genhead()

	var keys []string
	for key := range t.typeInfo {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		str += t.genSetFunc(key, t.typeInfo[key])
		str += t.genSetDirtyFunc(key)
		str += t.genGetFunc(key, t.typeInfo[key])
	}

	//str = str[:len(str)-1]
	str += fmt.Sprintf("type I%sDef interface {\n", t.name)
	str += t.interfaceStr
	str += "}\n"

	return str
}

// AddType 增加类型
func (t *Template) AddType(typename, typ string) {
	t.typeInfo[typename] = typ
	if strings.Contains(typ, "protoMsg") {
		t.needAddProto = true
	}
}

func (t *Template) genhead() string {
	var baseStr string

	baseStr += "package entitydef\n"
	baseStr += "\n"
	baseStr += "import \"zeus/iserver\"\n"
	if t.needAddProto {
		baseStr += "import \"protoMsg\"\n"
	}
	baseStr += "\n"
	baseStr += fmt.Sprintf("// %sDef 自动生成的属性包装代码\n", t.name)
	baseStr += fmt.Sprintf("type %sDef struct {\n", t.name)
	baseStr += "	ip iserver.IEntityProps\n"
	baseStr += "}\n"
	baseStr += "\n"
	baseStr += "// SetPropsSetter 设置接口\n"
	baseStr += fmt.Sprintf("func (p *%sDef) SetPropsSetter(ip iserver.IEntityProps) {\n", t.name)
	baseStr += "	p.ip = ip\n"
	baseStr += "}\n"
	baseStr += "\n"

	return baseStr
}

func (t *Template) genSetFunc(typename, typ string) string {
	var baseStr string
	baseStr += fmt.Sprintf("// Set%s 设置 %s\n", typename, typename)
	baseStr += fmt.Sprintf("func (p *%sDef) Set%s(v %s) {\n", t.name, typename, typ)
	baseStr += fmt.Sprintf("	p.ip.SetProp(\"%s\", v)\n", typename)
	baseStr += "}\n"
	baseStr += "\n"

	t.interfaceStr += fmt.Sprintf("	Set%s(v %s)\n", typename, typ)
	return baseStr
}

func (t *Template) genSetDirtyFunc(typename string) string {
	var baseStr string
	baseStr += fmt.Sprintf("// Set%sDirty 设置%s被修改\n", typename, typename)
	baseStr += fmt.Sprintf("func (p *%sDef) Set%sDirty() {\n", t.name, typename)
	baseStr += fmt.Sprintf("	p.ip.PropDirty(\"%s\")\n", typename)
	baseStr += "}\n"
	baseStr += "\n"

	t.interfaceStr += fmt.Sprintf("	Set%sDirty()\n", typename)
	return baseStr
}

func (t *Template) genGetFunc(typename, typ string) string {
	var baseStr string
	baseStr += fmt.Sprintf("// Get%s 获取 %s\n", typename, typename)
	baseStr += fmt.Sprintf("func (p *%sDef) Get%s() %s {\n", t.name, typename, typ)
	baseStr += fmt.Sprintf("	v := p.ip.GetProp(\"%s\")\n", typename)
	baseStr += "	if v == nil {\n"

	switch typ {
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		baseStr += fmt.Sprintf("		return %s(0)\n", typ)
	case "string":
		baseStr += "		return \"\"\n"
	case "bool":
		baseStr += "		return false\n"
	default:
		baseStr += "		return nil\n"
	}

	baseStr += "	}\n"
	baseStr += "\n"
	baseStr += fmt.Sprintf("	return v.(%s)\n", typ)
	baseStr += "}\n"
	baseStr += "\n"

	t.interfaceStr += fmt.Sprintf("	Get%s() %s\n", typename, typ)
	return baseStr
}
