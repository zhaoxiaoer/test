package obd

import (
	"testing"
)

func TestBCommEncode(t *testing.T) {
	bcomm := NewBComm()
	frames, err := bcomm.Encode([]byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Errorf("bcomm.Encode error: %v", err)
	} else {
		if bsTostr(frames[0]) != "AA0003010203B3" {
			t.Errorf("bcomm.Encode error")
		}
	}
}

func TestBCommDecode(t *testing.T) {
	bcomm := NewBComm()
	datas, err := bcomm.Decode([]byte{0xAA, 0x00, 0x03, 0x01, 0x02, 0x03, 0xB3})
	if err != nil {
		t.Errorf("bcomm.Decode error: %v", err)
	}
	if bsTostr(datas[0]) != "010203" {
		t.Errorf("bcomm.Decode error")
	}

	datas, err = bcomm.Decode([]byte{0xAB, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0xB0, 0xAC, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x03, 0xB4})
	if err != nil {
		t.Errorf("bcomm.Decode error: %v", err)
	}
	if bsTostr(datas[0]) != "03" {
		t.Errorf("bcomm.Decode error")
	}
}

func TestRecvAck(t *testing.T) {
	bcomm := NewBComm()

	go func() {
		for i := 0; i < 100000; i++ {
			bcomm.Decode([]byte{0xAD, 0x00, 0x04, 0x00, 0x00, 0x00, 0x03, 0xB4})
		}
	}()

	err := bcomm.recvACK(3)
	if err != nil {
		t.Errorf("bcomm.recvACK error")
	}
}
