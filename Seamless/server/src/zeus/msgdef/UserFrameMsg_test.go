package msgdef

import (
	"reflect"
	"testing"
)

func TestClientFrameMsgData(t *testing.T) {
	cfm := genClientFrameMsgData()

	data := make([]byte, cfm.Size())
	i, err := cfm.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != cfm.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, cfm.Size())
	}

	cfm2 := new(ClientFrameMsgData)
	if err = cfm2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(cfm, cfm2) {
		t.Fatalf("Before (%v), After (%v)", cfm, cfm2)
	}
}

func TestServerFrameMsg(t *testing.T) {
	sfm := genServerFrameMsg()

	data := make([]byte, sfm.Size())
	i, err := sfm.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != sfm.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, sfm.Size())
	}

	sfm2 := new(ServerFrameMsg)
	if err = sfm2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(sfm, sfm2) {
		t.Fatalf("Before (%v), After (%v)", sfm, sfm2)
	}
}

func TestFramesMsg(t *testing.T) {
	fm := genFramesMsg()

	data := make([]byte, fm.Size())
	i, err := fm.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != fm.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, fm.Size())
	}

	fm2 := new(FramesMsg)
	if err = fm2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(fm, fm2) {
		t.Fatalf("Before (%v), After (%v)", fm, fm2)
	}
}

func genClientFrameMsgData() *ClientFrameMsgData {
	cfm := new(ClientFrameMsgData)
	cfm.UID = 8
	cfm.Data = []byte{1, 2, 3, 4}
	return cfm
}

func genServerFrameMsg() *ServerFrameMsg {
	sfm := new(ServerFrameMsg)
	sfm.FrameID = 1024
	sfm.Msgs = make([]ClientFrameMsgData, 0, 1)
	sfm.Msgs = append(sfm.Msgs, *genClientFrameMsgData())
	return sfm
}

func genFramesMsg() *FramesMsg {
	fm := new(FramesMsg)
	//fm.Typ = 2
	fm.Frames = make([]ServerFrameMsg, 0, 1)
	fm.Frames = append(fm.Frames, *genServerFrameMsg())
	return fm
}
