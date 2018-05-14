package msgdef

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEnterAOI(t *testing.T) {
	msg := new(EnterAOI)
	msg.EntityType = "player"
	msg.EntityID = 10
	// msg.Pos = linmath.Vector3{
	// 	X: 10,
	// 	Y: 3.14,
	// 	Z: -10,
	// }
	// msg.Rota = linmath.Vector3{
	// 	X: 1024,
	// 	Y: 33,
	// 	Z: -88,
	// }
	msg.PropNum = 1
	msg.Properties = []byte("Hello World!")

	fmt.Println(msg)

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(EnterAOI)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestLeaveAOI(t *testing.T) {
	msg := new(LeaveAOI)
	msg.EntityID = 10

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(LeaveAOI)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestEnterSpace(t *testing.T) {
	msg := new(EnterCell)
	msg.SpaceID = 10

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(EnterCell)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestLeaveSpace(t *testing.T) {
	msg := new(LeaveCell)

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(LeaveCell)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

// func TestAOIPosChange(t *testing.T) {
// 	msg := NewAOIPosChangeMsg()

// 	msg.AddData([]byte("hello"))

// 	fmt.Println("msg1", msg.Num, string(msg.Data))

// 	data := make([]byte, msg.Size())
// 	i, err := msg.MarshalTo(data)
// 	fmt.Println(data)
// 	if err != nil {
// 		t.Fatalf("MarshalTo ERROR: %s", err)
// 	}
// 	if i != msg.Size() {
// 		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
// 	}

// 	msg2 := NewAOIPosChangeMsg()
// 	if err = msg2.Unmarshal(data); err != nil {
// 		t.Fatalf("Unmarshal ERROR: %s", err)
// 	}

// 	if !reflect.DeepEqual(msg, msg2) {
// 		t.Fatalf("Before (%v), After (%v)", msg, msg2)
// 	}

// 	fmt.Println("msg2", msg2.Num, string(msg2.Data))
// }
