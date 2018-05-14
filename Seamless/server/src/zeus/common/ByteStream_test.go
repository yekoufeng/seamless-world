package common

import (
	"fmt"
	"reflect"
	"testing"
)

type simStruct struct {
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
	B   byte
	BS  []byte
	Str string
	F32 float32
	F64 float64
}

func (s simStruct) Size() int {
	return 32 + len(s.BS) + len(s.Str)
}

func TestByteStreamMarshal(t *testing.T) {
	s := simStruct{}
	s.U8 = 7
	s.U16 = 16
	s.U32 = 55
	s.U64 = 100
	s.B = 10
	s.BS = []byte("test")
	s.Str = "hellow"
	s.F32 = 3.1415
	s.F64 = 3.1415926

	fmt.Println(s.Size())

	data := make([]byte, s.Size())

	bw := NewByteStream(data)

	fmt.Println(CalcSize(&s))

	bw.Marshal(&s)

	s1 := simStruct{}
	br := NewByteStream(data)

	err := br.Unmarshal(&s1)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, s1) {
		t.Fatalf("Before (%v), After (%v)", s, s1)
	}
}

func TestSimple(t *testing.T) {
	data := make([]byte, 4+8+1+12*4)
	bs := NewByteStream(data)
	fmt.Println(bs.data)
	if err := bs.WriteUInt64(100034); err != nil {
		fmt.Println(err)
	}
	fmt.Println(bs.data)
	if err := bs.WriteByte(1); err != nil {
		fmt.Println(err)
	}
	fmt.Println(bs.data)
	if err := bs.WriteInt32(32); err != nil {
		fmt.Println(err)
	}
	fmt.Println(bs.data)
}
