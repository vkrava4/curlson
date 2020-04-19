package app

import "github.com/Sirupsen/logrus"

type Configuration struct {
	Logs     *LogConfiguration
	Template *TemplateConfiguration
}

type TemplateConfiguration struct {
	Enabled bool
	Path    string
	Size    int
}

type LogConfiguration struct {
	Enabled        bool
	VerboseEnabled bool

	Log *logrus.Logger
}
