package msgdef

import (
	"reflect"
	"testing"
)

func TestClientVertifyReq(t *testing.T) {
	cv := new(ClientVertifyReq)
	cv.Source = 1
	cv.UID = 100001
	cv.Token = "38275d79848ce36d06e91bee644683ba"

	data := make([]byte, cv.Size())
	i, err := cv.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != cv.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, cv.Size())
	}

	cv2 := new(ClientVertifyReq)
	if err = cv2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(cv, cv2) {
		t.Fatalf("Before (%v), After (%v)", cv, cv2)
	}
}

func TestClientVertifySucceedRet(t *testing.T) {
	cvsr := new(ClientVertifySucceedRet)
	cvsr.Source = 1
	cvsr.UID = 100001
	cvsr.SourceID = 10
	cvsr.Type = 3

	data := make([]byte, cvsr.Size())
	i, err := cvsr.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != cvsr.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, cvsr.Size())
	}

	cvsr2 := new(ClientVertifySucceedRet)
	if err = cvsr2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(cvsr, cvsr2) {
		t.Fatalf("Before (%v), After (%v)", cvsr, cvsr2)
	}
}

func TestClientVertifyFailedRet(t *testing.T) {
	msg := new(ClientVertifyFailedRet)

	data := make([]byte, msg.Size())
	i, err := msg.MarshalTo(data)
	if err != nil {
		t.Fatalf("MarshalTo ERROR: %s", err)
	}
	if i != msg.Size() {
		t.Fatalf("MarshaTo return Size(%d), want (%d)", i, msg.Size())
	}

	msg2 := new(ClientVertifyFailedRet)
	if err = msg2.Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal ERROR: %s", err)
	}

	if !reflect.DeepEqual(msg, msg2) {
		t.Fatalf("Before (%v), After (%v)", msg, msg2)
	}
}
