package serializer

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSerializer(t *testing.T) {
	data, _ := Serialize(float32(3.14))
	fmt.Printf("%x\n", data)
	fmt.Println(data)

	ret, err := UnSerialize(data)
	fmt.Println(ret, err)
	for _, v := range ret {
		fmt.Println(reflect.TypeOf(v), v)
	}
}

func TestSerializerNew(t *testing.T) {
	// data := SerializeNew(uint8(8), uint16(1024), uint32(32), uint64(64),
	// 	int8(-8), int16(-16), int32(-32), int64(-64),
	// 	float32(3.1415), float64(-3.1415926), "helloworld", true)
	data := SerializeNew()
	fmt.Println(data)

	ret := UnSerializeNew(data)
	for _, v := range ret {
		fmt.Println(reflect.TypeOf(v), v)
	}

	// msg := &msgdef.EnterAOI{
	// 	EntityID:   10001,
	// 	EntityType: "player",
	// 	Pos:        linmath.Vector3{1024, 1024, -1024},
	// 	PropNum:    1,
	// 	Properties: []byte("hello"),
	// }
	// data = SerializeNew(msg)
	// fmt.Println(data)
	// ret = UnSerializeNew(data)
	// fmt.Println(ret)
}
