package ot

import (
	"log"
	"sync"
)

type OTBufferManager struct {
	OTBuffers OTBufferMap
}

var (
	otBuffermanagerInstance *OTBufferManager
	once                    sync.Once
)

func initializeBufferManager() {
	once.Do(func() {
		otBuffermanagerInstance = &OTBufferManager{
			OTBuffers: make(OTBufferMap),
		}
	})
}

func GetOTBufferManager() *OTBufferManager {
	initializeBufferManager()
	return otBuffermanagerInstance
}

func (otbm *OTBufferManager) GetOTBuffer(documentUUID string) *OTBuffer {
	otBuffer, ok := otbm.OTBuffers[documentUUID]
	if !ok {
		otBuffer = NewOTBuffer(documentUUID)
		otbm.OTBuffers[documentUUID] = otBuffer
	}
	return otBuffer
}

func (otbm *OTBufferManager) ProcessTransformations() {
	log.Print("processing transformations...")
	for {
		for uuid := range otbm.OTBuffers {
			otbm.OTBuffers[uuid].process()
		}
	}
}
