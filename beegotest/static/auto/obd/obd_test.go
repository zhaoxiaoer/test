package obd

import (
	"testing"
)

func TestOBDInit(t *testing.T) {
	obd := NewOBD()
	err := obd.Init()
	if err != nil {
		t.Errorf("obd.Init error: %s", err)
	}
	err = obd.Init()
	if err == nil {
		t.Errorf("obd.Init error")
	}
}

func TestOBDUninit(t *testing.T) {
	obd := NewOBD()
	obd.Init()
	err := obd.Uninit()
	if err != nil {
		t.Error("obd.Uninit error: %s", err)
	}
	err = obd.Uninit()
	if err == nil {
		t.Errorf("obd.Uninit error")
	}
}

func TestOBDWrite(t *testing.T) {
	obd := NewOBD()
	obd.Init()
	err := obd.Write([]byte("123"))
	if err != nil {
		t.Errorf("obd.Write error: %s\n", err)
	}
}

func TestOBDIsConnected(t *testing.T) {
	obd := NewOBD()
	obd.Init()
	_, err := obd.IsConnected()
	if err != nil {
		t.Errorf("obd.IsConnected error: %s\n", err)
	}
}
