package obd

import (
	"testing"
)

func TestTCPServerInit(t *testing.T) {
	ts := NewTCPServer()
	err := ts.Init()
	if err != nil {
		t.Errorf("tcpserver.Init error: %s", err)
	}
	ts.Uninit()
}

func TestTCPServerUninit(t *testing.T) {
	ts := NewTCPServer()
	ts.Init()
	err := ts.Uninit()
	if err != nil {
		t.Errorf("tcpserver.Uninit error: %s", err)
	}
}

func TestTCPServerWrite(t *testing.T) {
	ts := NewTCPServer()
	ts.Init()
	if err := ts.Write([]byte("123")); err != nil {
		t.Errorf("tcpserver.Write error: %s", err)
	}
	ts.Uninit()
}

func TestTCPServerIsConnected(t *testing.T) {
	ts := NewTCPServer()
	ts.Init()
	if _, err := ts.IsConnected(); err != nil {
		t.Errorf("tcpserver.IsConnected error: %s", err)
	}
	ts.Uninit()
}
