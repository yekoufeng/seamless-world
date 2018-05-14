package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"zeus/iserver"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
)

// Defs 所有的实体定义
type Defs struct {
	defs map[string]*Def
}

var inst *Defs

// GetDefs 获取实体定义信息
func GetDefs() *Defs {
	return inst
}

func initDefs() {
	inst = &Defs{
		defs: make(map[string]*Def),
	}

	inst.Init()
}

// Init 初始化文件定义结构
func (defs *Defs) Init() {
	err := filepath.Walk("../res/entitydef", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}

		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				log.Error("read entity def file error ", err)
				return nil
			}

			jsonInfo := make(map[string]interface{})
			err = json.Unmarshal(raw, &jsonInfo)
			if err != nil {
				log.Error("parse entity def file error ", err)
				return nil
			}

			def := newDef()
			if err := def.fill(jsonInfo); err != nil {
				log.Error("fill def error ", err)
				return nil
			}

			defs.defs[jsonInfo["name"].(string)] = def
		}

		return nil
	})

	if err != nil {
		log.Error("walk entity def file error")
	}
}

// GetDef 获取一个entity定义
func (defs *Defs) GetDef(name string) *Def {
	d, ok := defs.defs[name]
	if !ok {
		return nil
	}

	return d
}

////////////////////////////////////////////////////////////////

// Def 独立的实体定义
type Def struct {
	Name        string
	Props       map[string]*PropDef
	ClientProps map[string]uint8
}

func newDef() *Def {
	return &Def{
		Props:       make(map[string]*PropDef),
		ClientProps: make(map[string]uint8),
	}
}

func (def *Def) fill(jsonInfo map[string]interface{}) error {

	def.Name = jsonInfo["name"].(string)

	jsonProps := jsonInfo["props"].(map[string]interface{})
	jsonServers := jsonInfo["server"].(map[string]interface{})
	jsonClient, ok := jsonInfo["client"].(map[string]interface{})["props"].([]interface{})
	if !ok {
		log.Error("读取客户端关注的属性失败!")
	}
	jsonMRoleMap, ok := jsonInfo["mrole"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("读取失败")
	}
	jsonMRole, ok := jsonMRoleMap["props"].([]interface{})
	if !ok {
		log.Error("读取主角关注的属性失败!")
	}

	for propName, propInfo := range jsonProps {
		jsonProp := propInfo.(map[string]interface{})

		prop := newPropDef()
		def.Props[propName] = prop
		prop.Name = propName
		prop.Desc = jsonProp["desc"].(string)
		prop.Type = def.getTypeByStr(jsonProp["type"].(string))
		prop.TypeName = jsonProp["type"].(string)
		if strings.Contains(prop.TypeName, "protoMsg") {
			prop.TypeName = prop.TypeName[10:]
		}
		prop.Persistence = true
		if persistence, ok := jsonProp["save"].(string); ok {
			if persistence == "0" {
				prop.Persistence = false
			}
		}
		prop.Sync = true
		if sync, ok := jsonProp["sync"].(string); ok {
			if sync == "0" {
				prop.Sync = false
			}
		}

		for srvTypeName, srvInfo := range jsonServers {

			srvPropList, ok := srvInfo.(map[string]interface{})["props"].([]interface{})
			if !ok {
				log.Error("分析服务器关心的属性列表失败", srvTypeName, srvInfo)
				continue
			}

			isInterest := false
			for _, srvPropName := range srvPropList {
				if srvPropName.(string) == propName {
					isInterest = true
					break
				}
			}

			if isInterest {
				srvType, _ := strconv.Atoi(srvTypeName)
				prop.InterestSrvs = append(prop.InterestSrvs, uint8(srvType))
			}
		}

		prop.IsClientInterest = false
		for _, clientProp := range jsonClient {
			if propName == clientProp.(string) {
				prop.IsClientInterest = true
				break
			}
		}

		prop.IsMRoleInterest = false
		for _, mroleProp := range jsonMRole {
			if propName == mroleProp.(string) {
				prop.IsMRoleInterest = true
				prop.InterestSrvs = append(prop.InterestSrvs, iserver.ServerTypeGateway)
				break
			}
		}
	}

	return nil
}

func (def *Def) getTypeByStr(st string) reflect.Type {

	var t reflect.Type

	switch st {
	case "int8":
		t = reflect.TypeOf(int8(0))
	case "int16":
		t = reflect.TypeOf(int16(0))
	case "int32":
		t = reflect.TypeOf(int32(0))
	case "int64":
		t = reflect.TypeOf(int64(0))
	case "byte":
		t = reflect.TypeOf(byte(0))
	case "uint8":
		t = reflect.TypeOf(uint8(0))
	case "uint16":
		t = reflect.TypeOf(uint16(0))
	case "uint32":
		t = reflect.TypeOf(uint32(0))
	case "uint64":
		t = reflect.TypeOf(uint64(0))
	case "float32":
		t = reflect.TypeOf(float32(0))
	case "float64":
		t = reflect.TypeOf(float64(0))
	case "string":
		t = reflect.TypeOf("")
	case "bool":
		t = reflect.TypeOf(false)
	default:
		return nil
	}

	return t
}

////////////////////////////////////////////////////////////////

// PropDef 字段定义
type PropDef struct {
	Name             string
	Desc             string
	Type             reflect.Type
	TypeName         string
	InterestSrvs     []uint8
	IsClientInterest bool
	IsMRoleInterest  bool
	Persistence      bool //是否需要持久化, 默认true
	Sync             bool //是否需要同步至其他服务器, 默认true
}

func newPropDef() *PropDef {
	return &PropDef{
		InterestSrvs: make([]uint8, 0, 10),
	}
}

// CreateInst 创建该属性的一个实例
func (pd *PropDef) CreateInst() (interface{}, error) {

	if pd.Type != nil {
		return reflect.New(pd.Type), nil
	}

	msgID, err := msgdef.GetMsgDef().GetMsgIDByName(pd.TypeName)
	if err != nil {
		return nil, err
	}

	_, msgContent, err := msgdef.GetMsgDef().GetMsgInfo(msgID)
	return msgContent, err
}

// IsValidValue 当前值是否能设置到该属性上
func (pd *PropDef) IsValidValue(value interface{}) bool {

	if value == nil {
		return true
	}

	if pd.Type != nil {
		return pd.Type == reflect.TypeOf(value)
	}

	return strings.Contains(reflect.TypeOf(value).String(), pd.TypeName)
}
