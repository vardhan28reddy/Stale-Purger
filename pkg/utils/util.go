package utils

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

var StalePodStates = map[v1.PodPhase]bool{
	v1.PodUnknown: true,
	v1.PodFailed:  true,
}

var FatalFunc = func(errInfo string, err error, logger *logrus.Entry) {
	if err != nil {
		logger.Fatalf("%s: %+v", errInfo, err)
	}
}
