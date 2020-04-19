package util

import "github.com/Sirupsen/logrus"

type AppConfiguration struct {
	logs     *LogConfiguration
	template *TemplateConfiguration
}

type TemplateConfiguration struct {
	enabled bool
	path    string
	size    int
}

type LogConfiguration struct {
	enabled        bool
	verboseEnabled bool

	log *logrus.Logger
}
