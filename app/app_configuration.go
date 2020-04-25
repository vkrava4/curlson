package app

import (
	"github.com/Sirupsen/logrus"
	"os"
)

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
	Enabled bool
	Persist bool
	Verbose bool

	Log     *logrus.Logger
	Logfile *os.File
}
