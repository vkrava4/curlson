package util

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"os"
	"strings"
)

var (
	fileBuffer            = 32 * 1024
	filesMode             = os.FileMode(0666)
	filesEndLineDelimiter = byte('\n')
)

func CountLinesForFile(file *os.File) int {
	var reader = bufio.NewReader(file)
	var buffer = make([]byte, fileBuffer)
	var count = 0
	var line = []byte{'\n'}

	for {
		var c, err = reader.Read(buffer)
		count += bytes.Count(buffer[:c], line)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			return -1
		}
	}
}

// Counts the number of lines in a given file.
// In case property 'templateEnabled' is false (i.e templating is disabled) either there is an error occurred while reading a file
// the result '-1' will be returned
func CountLines(filePath string, templateEnabled bool) int {
	if !templateEnabled {
		return -1
	}

	var file, _ = os.OpenFile(filePath, os.O_RDONLY, filesMode)
	defer file.Close()

	var reader = bufio.NewReader(file)
	var buffer = make([]byte, fileBuffer)
	var count = 0
	var line = []byte{'\n'}

	for {
		var c, err = reader.Read(buffer)
		count += bytes.Count(buffer[:c], line)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
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

	var file, _ = os.OpenFile(templateFile, os.O_RDONLY, filesMode)
	defer file.Close()
	var reader = bufio.NewReader(file)

	for {
		counter++
		if counter == lineNumber {
			var line, _ = reader.ReadString(filesEndLineDelimiter)
			return strings.TrimSuffix(line, string(filesEndLineDelimiter))
		} else {
			var _, _ = reader.ReadString(filesEndLineDelimiter)
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
