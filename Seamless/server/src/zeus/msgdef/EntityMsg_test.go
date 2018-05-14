package msgdef

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClientTransport(t *testing.T) {
	ct := new(ClientTransport)
	ct.SrvType = 1
	ct.MsgFlag = 2
	ct.MsgContent = []byte("Test")

	data := make([]byte, ct.Size())
	i, err := ct.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != ct.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, ct.Size())
	}

	ct2 := new(ClientTransport)
	if err = ct2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(ct, ct2) {
		t.Fatalf("Before (%v), After (%v)", ct, ct2)
	}
}

func TestCreateEntityReq(t *testing.T) {
	ce := new(CreateEntityReq)
	ce.EntityType = "Test"
	ce.EntityID = 1
	ce.SpaceID = 5000
	ce.DBID = 1024
	// ce.InitParam = "5"
	ce.SrcSrvType = 2
	ce.SrcSrvID = 1000
	ce.CallbackID = 123

	data := make([]byte, ce.Size())
	i, err := ce.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	fmt.Println(data)
	if i != ce.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, ce.Size())
	}

	ce2 := new(CreateEntityReq)
	if err = ce2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}
	fmt.Println(ce2)
	if !reflect.DeepEqual(ce, ce2) {
		t.Fatalf("Before (%v), After (%v)", ce, ce2)
	}
}

func TestCreateEntityRet(t *testing.T) {
	cer := new(CreateEntityRet)
	cer.SrvType = 3
	cer.EntityID = 1001
	cer.CallbackID = 234
	cer.ErrorStr = "Testing"

	data := make([]byte, cer.Size())
	i, err := cer.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != cer.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, cer.Size())
	}

	cer2 := new(CreateEntityRet)
	if err = cer2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(cer, cer2) {
		t.Fatalf("Before (%v), After (%v)", cer, cer2)
	}
}

func TestDestroyEntityReq(t *testing.T) {
	de := new(DestroyEntityReq)
	de.EntityID = 1100
	de.SrcSrvType = 8
	de.SrcSrvID = 2000
	de.CallbackID = 345

	data := make([]byte, de.Size())
	i, err := de.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != de.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, de.Size())
	}

	de2 := new(DestroyEntityReq)
	if err = de2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(de, de2) {
		t.Fatalf("Before (%v), After (%v)", de, de2)
	}
}

func TestDestroyEntityRet(t *testing.T) {
	der := new(DestroyEntityRet)
	der.SrvType = 3
	der.EntityID = 3000
	der.CallbackID = 18
	der.ErrorStr = "Testing"

	data := make([]byte, der.Size())
	i, err := der.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != der.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, der.Size())
	}

	der2 := new(DestroyEntityRet)
	if err = der2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(der, der2) {
		t.Fatalf("Before (%v), After (%v)", der, der2)
	}
}

func TestEntityMsgTransport(t *testing.T) {
	msg := new(EntityMsgTransport)
	msg.SrvType = 5
	msg.EntityID = 1024
	msg.CellID = 50
	msg.MsgContent = []byte("hello world")

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(EntityMsgTransport)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestEntityVarData(t *testing.T) {
	msg := new(EntityVarData)
	msg.Identifier = "test"
	msg.Variant = []byte("hello world")

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(EntityVarData)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}

func TestEntityMsgChange(t *testing.T) {
	evd := new(EntityVarData)
	evd.Identifier = "test"
	evd.Variant = []byte{1, 2, 3, 4}

	emc := new(EntityMsgChange)
	emc.VarData = append(emc.VarData, *evd)

	data := make([]byte, emc.Size())
	i, err := emc.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != emc.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, emc.Size())
	}

	emc2 := new(EntityMsgChange)
	if err = emc2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(emc, emc2) {
		t.Fatalf("Before (%v), After (%v)", emc, emc2)
	}
}

func TestRPCMsg(t *testing.T) {
	rpc := new(RPCMsg)
	rpc.ServerType = 4
	rpc.SrcEntityID = 1024
	rpc.MethodName = "test"
	rpc.Data = []byte("Hello World!")

	data := make([]byte, rpc.Size())
	i, err := rpc.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != rpc.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, rpc.Size())
	}

	rpc2 := new(RPCMsg)
	if err = rpc2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(rpc, rpc2) {
		t.Fatalf("Before (%v), After (%v)", rpc, rpc2)
	}
}
