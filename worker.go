package main

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func pinWorker(workerID int, ptrQueue chan *MetaPointer) {
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
			l.WithFields(logrus.Fields{
				"metadataProtocol": ptr.Protocol,
				"metadataPointer":  ptr.Pointer,
			}).Info("pinning metaPointer")
		}
	}
}
