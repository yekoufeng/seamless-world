package msgdef

import (
	"reflect"
	"testing"
)

func TestSrvMsgTransport(t *testing.T) {
	msg := new(SrvMsgTransport)
	msg.SpaceID = 3000
	msg.MsgContent = []byte("test")

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(SrvMsgTransport)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}
