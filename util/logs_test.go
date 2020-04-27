package util

import (
	"github.com/sirupsen/logrus"
	"github.com/vkrava4/curlson/app"
	"os"
	"strings"
	"testing"
)

func TestSetupLogs(t *testing.T) {
	var givenConfig = &app.LogConfiguration{
		Enabled: false,
		Persist: false,
		Verbose: false,
		Log:     logrus.New(),
		Logfile: nil,
	}

	var errSetup = SetupLogs(givenConfig)
	defer cleanup(givenConfig.Logfile.Name())

	if errSetup != nil {
		t.Errorf("Result is incorrect. An error is not expected: %v", givenConfig)
	}

	if givenConfig.Logfile == nil || !strings.Contains(givenConfig.Logfile.Name(), "current-execution-") ||
		!givenConfig.Enabled || givenConfig.Log.Out == nil {
		t.Errorf("Result is incorrect: %v", givenConfig)
	}

}

func TestShutdownLogsAndRename(t *testing.T) {
	var givenConfig = &app.LogConfiguration{
		Enabled: false,
		Persist: true,
		Verbose: false,
		Log:     logrus.New(),
		Logfile: nil,
	}

	var errSetup = SetupLogs(givenConfig)
	defer cleanup(givenConfig.Logfile.Name())

	if errSetup != nil {
		t.Errorf("Result is incorrect. An error is not expected: %v", givenConfig)
	}

	var logfileNameToBeRenamed = givenConfig.Logfile.Name()
	ShutdownLogs(givenConfig)

	if givenConfig.Logfile != nil || givenConfig.Enabled || givenConfig.Log != nil {
		t.Errorf("Result is incorrect: %v", givenConfig)
	}

	var _, infoErr1 = os.Stat(logfileNameToBeRenamed)
	if !os.IsNotExist(infoErr1) {
		t.Errorf("File should be renamed: %s", logfileNameToBeRenamed)
	}

	var info2, infoErr2 = os.Stat(strings.Replace(logfileNameToBeRenamed, "current-execution-", "execution-", -1))
	if infoErr2 != nil || info2.IsDir() {
		t.Errorf("File should be renamed: %s", logfileNameToBeRenamed)
	}
}

func TestShutdownLogsAndRemove(t *testing.T) {
	var givenConfig = &app.LogConfiguration{
		Enabled: false,
		Persist: false,
		Verbose: false,
		Log:     logrus.New(),
		Logfile: nil,
	}

	var errSetup = SetupLogs(givenConfig)
	defer cleanup(givenConfig.Logfile.Name())

	if errSetup != nil {
		t.Errorf("TestShutdownLogs result is incorrect. An error is not expected: %v", givenConfig)
	}

	var logfileNameToBeRemoved = givenConfig.Logfile.Name()
	ShutdownLogs(givenConfig)

	if givenConfig.Logfile != nil || givenConfig.Enabled || givenConfig.Log != nil {
		t.Errorf("TestShutdownLogsAndRemove is incorrect: %v", givenConfig)
	}

	var _, infoErr = os.Stat(logfileNameToBeRemoved)
	if !os.IsNotExist(infoErr) {
		t.Errorf("File should be removed: %s", logfileNameToBeRemoved)
	}
}

func TestWarnLogIfEnabled(t *testing.T) {
	var givenConfig = &app.LogConfiguration{
		Enabled: true,
		Log:     logrus.New(),
	}

	WarnLog("Test", givenConfig)
}

func TestErrorLogIfEnabled(t *testing.T) {
	var givenConfig = &app.LogConfiguration{
		Enabled: true,
		Log:     logrus.New(),
	}

	ErrorLog("Test", givenConfig)
}

func TestInfoLogIfEnabled(t *testing.T) {
	var givenConfig = &app.LogConfiguration{
		Enabled: true,
		Log:     logrus.New(),
	}

	InfoLog("Test", givenConfig)
}

func cleanup(file string) {
	_ = os.Remove(file)
}
