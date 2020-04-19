package util

import "github.com/Sirupsen/logrus"

type AppConfiguration struct {
	templatingEnabled bool
	loggingEnabled    bool
	verboseEnabled    bool

	log *logrus.Logger
}
