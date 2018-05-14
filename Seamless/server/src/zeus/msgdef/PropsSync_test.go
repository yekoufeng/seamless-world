package msgdef

import (
	"reflect"
	"testing"
)

func TestPropsSync(t *testing.T) {
	msg := new(PropsSync)
	msg.Num = 1024
	msg.Data = []byte("hello world!")

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(PropsSync)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestPropsSyncClient(t *testing.T) {
	msg := new(PropsSyncClient)
	msg.EntityID = 1024
	msg.Num = 10
	msg.Data = []byte("hello world!")

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(PropsSyncClient)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}
