package ot

import (
	"sync"
)

type OTBufferManager struct {
	OTBuffers OTBufferMap
}

var (
	otBuffermanagerInstance *OTBufferManager
	once                    sync.Once
)

// initialize the buffer manager only once
func initializeBufferManager() {
	once.Do(func() {
		otBuffermanagerInstance = &OTBufferManager{
			OTBuffers: make(OTBufferMap),
		}
	})
}

// returns the current OTBufferManager and creates on if it isn't created
func GetOTBufferManager() *OTBufferManager {
	initializeBufferManager()
	return otBuffermanagerInstance
}

// returns the OTBuffer given the documentUUID and creates one if it doesn't exists
func (otbm *OTBufferManager) GetOTBuffer(documentUUID string) *OTBuffer {
	otBuffer, ok := otbm.OTBuffers[documentUUID]
	if !ok {
		otBuffer = NewOTBuffer(documentUUID)
		otbm.OTBuffers[documentUUID] = otBuffer
	}
	return otBuffer
}

// The loop that process pending transformation for all the documents
func (otbm *OTBufferManager) ProcessTransformations() {
	for {
		for uuid := range otbm.OTBuffers {
			go otbm.OTBuffers[uuid].Process()
		}
	}
}
