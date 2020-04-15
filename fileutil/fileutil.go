package fileutil

import (
	"bufio"
	"bytes"
	"curlson/logutil"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"strings"
)

var (
	bufferSize       = 32 * 1024
	defaultMode      = os.FileMode(0666)
	endLineDelimiter = byte('\n')
)

// Counts the number of lines in a given file.
// In case property 'templateEnabled' is false (i.e templating is disabled) either there is an error occurred while reading a file
// the result '-1' will be returned
func CountLines(templateFile *string, templateEnabled *bool, log *logrus.Logger, loggingSupported *bool) int {
	if !*templateEnabled {
		logutil.InfoLog("Skipping templating stage. It is disabled", log, loggingSupported)
		return -1
	}

	logutil.InfoLog(fmt.Sprintf("Determinig the number of lines in template file: %s", *templateFile), log, loggingSupported)
	var file, _ = os.OpenFile(*templateFile, os.O_RDONLY, defaultMode)
	defer file.Close()

	var reader = bufio.NewReader(file)

	var buffer = make([]byte, bufferSize)
	var count = 0
	var line = []byte{'\n'}

	for {
		var c, err = reader.Read(buffer)
		count += bytes.Count(buffer[:c], line)

		switch {
		case err == io.EOF:
			logutil.InfoLog(fmt.Sprintf("Reached the end of template file. Total lines: %d", count), log, loggingSupported)
			return count

		case err != nil:
			logutil.ErrorLog(fmt.Sprintf("An error occured while reading a file: '%s'. Error: %s", *templateFile, err.Error()), log, loggingSupported)
			return -1
		}
	}
}

func ReadRandomLine(templateFile string, linesCount int) (int, string) {
	var lineNum = rand.Intn(linesCount)
	return lineNum, ReadLine(templateFile, lineNum)
}

func ReadLine(templateFile string, lineNumber int) string {
	var counter = -1

	var file, _ = os.OpenFile(templateFile, os.O_RDONLY, defaultMode)
	defer file.Close()
	var reader = bufio.NewReader(file)

	for {
		counter++
		if counter == lineNumber {
			var line, _ = reader.ReadString(endLineDelimiter)
			return strings.TrimSuffix(line, string(endLineDelimiter))
		} else {
			var _, _ = reader.ReadString(endLineDelimiter)
		}
	}
}

func FileExists(path string) bool {
	if _, isNotExistErr := os.Stat(path); !os.IsNotExist(isNotExistErr) {
		return true
	} else {
		return false
	}
}
