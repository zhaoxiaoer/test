package obd

import (
	"fmt"
	"sync"
)

type OBD struct {
	sync.RWMutex
	hasInited bool

	Event chan string
}

func NewOBD() *OBD {
	obd := &OBD{}
	return obd
}

func (obd *OBD) Init() error {
	obd.Lock()
	defer obd.Unlock()

	if !obd.hasInited {
		obd.Event = make(chan string)
		// TODO

		obd.hasInited = true
		return nil
	} else {
		return fmt.Errorf("obd has been initialized!")
	}
}

func (obd *OBD) Uninit() error {
	obd.Lock()
	defer obd.Unlock()

	if obd.hasInited {
		// TODO

		close(obd.Event)
		obd.hasInited = false
		return nil
	} else {
		return fmt.Errorf("obd not initialized!")
	}
}

func (obd *OBD) Write(data []byte) error {
	obd.RLock()
	defer obd.RUnlock()

	if obd.hasInited {
		// TODO

		return nil
	} else {
		return fmt.Errorf("obd not initialized!")
	}
}

func (obd *OBD) IsConnected() (bool, error) {
	obd.RLock()
	defer obd.RUnlock()

	if obd.hasInited {
		// TODO

		return false, nil
	} else {
		return false, fmt.Errorf("obd not initialized!")
	}
}
