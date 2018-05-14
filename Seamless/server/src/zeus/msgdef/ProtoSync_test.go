package msgdef

import (
	"reflect"
	"testing"
)

func TestProtoSync(t *testing.T) {
	ps := new(ProtoSync)
	ps.Data = []byte("helloworld")

	data := make([]byte, ps.Size())
	i, err := ps.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != ps.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, ps.Size())
	}

	ps2 := new(ProtoSync)
	if err = ps2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(ps, ps2) {
		t.Fatalf("Before (%v), After (%v)", ps, ps2)
	}
}
