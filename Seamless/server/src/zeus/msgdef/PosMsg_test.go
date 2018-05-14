package msgdef

import (
	"testing"
)

func TestUserMove(t *testing.T) {
	/*
		msg := new(UserMove)
		msg.Pos = linmath.Vector3{
			X: 3.1415,
			Y: 0,
			Z: -3.1415,
		}

		data := make([]byte, msg.Size())
		i, err := msg.MarshalTo(data)
		fmt.Println(data)
		if err != nil {
			t.Fatalf("MarshalTo ERROR: %s", err)
		}
		if i != msg.Size() {
			t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
		}

		msg2 := new(UserMove)
		if err = msg2.Unmarshal(data); err != nil {
			t.Fatalf("Unmarshal ERROR: %s", err)
		}

		if !reflect.DeepEqual(msg, msg2) {
			t.Fatalf("Before (%v), After (%v)", msg, msg2)
		}
	*/
}
