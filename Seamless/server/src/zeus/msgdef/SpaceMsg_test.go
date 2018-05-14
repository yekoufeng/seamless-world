package msgdef

import (
	"reflect"
	"testing"
)

func TestEnterSpaceReq(t *testing.T) {
	msg := new(EnterCellReq)
	msg.SrvID = 1000
	msg.SpaceID = 2000
	msg.EntityType = "test"
	// msg.InitParam = "1"
	msg.EntityID = 3000
	msg.DBID = 4000
	msg.OldSrvID = 5000
	msg.OldSpaceID = 6000

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(EnterCellReq)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestLeaveSpaceReq(t *testing.T) {
	msg := new(LeaveCellReq)
	msg.EntityID = 3000

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(LeaveCellReq)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}
