package util

import (
	"github.com/fatih/color"
	"github.com/vkrava4/curlson/app"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var redColor = color.New(color.FgRed)
var yellowColor = color.New(color.FgYellow)

func SetupLogs(config *app.LogConfiguration) error {

	var defaultLogName = "current-execution-" + time.Now().Format(time.RFC3339) + ".log"
	var defaultAppDir, errAppDir = getApplicationDirectory()
	if errAppDir != nil {
		return errAppDir
	}

	var logfile, errCreateFile = os.OpenFile(filepath.Join(defaultAppDir, defaultLogName), os.O_WRONLY|os.O_CREATE, filesMode)
	if errCreateFile != nil {
		return errCreateFile
	}

	config.Log.SetOutput(logfile)
	config.Enabled = true
	config.Logfile = logfile

	return nil
}

func ShutdownLogs(config *app.LogConfiguration) {
	if config.Logfile == nil {
		return
	}

	var closeFileErr = config.Logfile.Close()
	if closeFileErr == nil {
		if config.Persist {
			var newLogfileName = strings.Replace(config.Logfile.Name(), "current-execution-", "execution-", -1)
			var _ = os.Rename(config.Logfile.Name(), newLogfileName)
		} else {
			var _ = os.Remove(config.Logfile.Name())
		}
	}
}

func InfoLog(s string, config *app.LogConfiguration) {
	if config.Enabled {
		config.Log.Info(s)
	}
}

func WarnLog(s string, config *app.LogConfiguration) {
	if config.Enabled {
		config.Log.Warn(s)
	}
}

func ErrorLog(s string, config *app.LogConfiguration) {
	if config.Enabled {
		config.Log.Error(s)
	}
}
