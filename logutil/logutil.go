package logutil

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"os"
	"strings"
	"time"
)

var defaultFolderName = ".curlson"

var redColor = color.New(color.FgRed)
var yellowColor = color.New(color.FgYellow)

func SetupLogs(log *logrus.Logger, persistLogs *bool, verbose *bool) (bool, *os.File) {
	var homeDir, homeDirErr = os.UserHomeDir()

	if homeDirErr == nil {
		var defaultFolderPath = homeDir + "/" + defaultFolderName
		if _, isNotExistErr := os.Stat(defaultFolderPath); os.IsNotExist(isNotExistErr) {
			var mkdirErr = os.Mkdir(defaultFolderPath, os.ModePerm)
			if mkdirErr != nil {
				fmt.Println(redColor.Sprintf("Application can not create directory: '%s'. Reason: %s", defaultFolderName, mkdirErr.Error()))
			}
		}

		var defaultLogFileName = "latest-execution-" + time.Now().Format(time.RFC3339) + ".log"
		var logfile, createFileErr = os.OpenFile(defaultFolderPath+string(os.PathSeparator)+defaultLogFileName, os.O_WRONLY|os.O_CREATE, 0666)

		if createFileErr == nil {
			log.SetOutput(logfile)
			var logsActionMessage = "auto cleaned"

			if *persistLogs {
				logsActionMessage = fmt.Sprintf("transformed to persistent file '%s'", strings.Replace(logfile.Name(), "latest-execution-", "execution-", -1))
			}

			if *verbose {
				fmt.Println(fmt.Sprintf("Created temporary log file '%s' which will be %s after execution", logfile.Name(), logsActionMessage))
			}
			return true, logfile
		} else if *verbose {
			fmt.Println(yellowColor.Sprintf("Application can not create temporary log file '%s'. Reason: %s", defaultFolderPath+string(os.PathSeparator)+defaultLogFileName, createFileErr.Error()))
		}

	} else if *verbose {
		fmt.Println(yellowColor.Sprintf("Application can not determine user's HOME directory. Reason: %s", homeDirErr.Error()))
	}

	return false, nil
}

func ShutdownLogs(logfile *os.File, persistLogs *bool, verbose *bool) {
	if logfile == nil {
		return
	}

	var closeFileErr = logfile.Close()
	if closeFileErr == nil {
		if *persistLogs {
			var newLogfileName = strings.Replace(logfile.Name(), "latest-execution-", "execution-", -1)
			var renameErr = os.Rename(logfile.Name(), newLogfileName)
			if renameErr != nil && *verbose {
				fmt.Println(yellowColor.Sprintf("Application can not rename temporary log file: '%s' with new name '%s'. Reason: %s", logfile.Name(), newLogfileName, renameErr.Error()))
			}
		} else {
			var removeErr = os.Remove(logfile.Name())
			if removeErr != nil && *verbose {
				fmt.Println(yellowColor.Sprintf("Application can not cleanup a log file: '%s'. Reason: %s", logfile.Name(), removeErr.Error()))
			}
		}
	} else if *verbose {
		fmt.Println(yellowColor.Sprintf("Application can not close a log file: '%s'. Reason: %s", logfile.Name(), closeFileErr.Error()))
	}
}

func InfoLog(s string, log *logrus.Logger, supported *bool) {
	if *supported {
		log.Info(s)
	}
}

func WarnLog(s string, log *logrus.Logger, supported *bool) {
	if *supported {
		log.Warn(s)
	}
}

func ErrorLog(s string, log *logrus.Logger, supported *bool) {
	if *supported {
		log.Error(s)
	}
}
