package main

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func pinningWorker(workerID int, ptrQueue chan *MetaPointer) {
	localLogger := logger.WithFields(logrus.Fields{
		"workerID": workerID,
		"action":   "pinWorker",
	})

	localLogger.Info("started")

	for {
		jobUUID := uuid.New()
		l := localLogger.WithFields(logrus.Fields{
			"jobUUID": jobUUID.String(),
		})

		select {
		case ptr := <-ptrQueue:
			ptrLogger := l.WithFields(logrus.Fields{
				"metadataProtocol": ptr.Protocol,
				"metadataPointer":  ptr.Pointer,
			})
			ptrLogger.Info("pinning metaPointer")

			err := pinningClient.Pin(ptr.Pointer)
			if err != nil {
				ptrLogger.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("error pinning")
			} else {
				ptrLogger.Info("pinned")
			}
		}
	}
}
