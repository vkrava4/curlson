package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryExist(t *testing.T) {
	var testDirName = "test"
	var absTestDirName, _ = filepath.Abs(testDirName)
	var _ = os.Mkdir(absTestDirName, os.ModePerm)
	defer os.Remove(absTestDirName)

	var exists = directoryExist(absTestDirName)

	if !exists {
		t.Error("Invalid result: directory should exist")
	}
}

func TestDirectoryDoesntExist(t *testing.T) {
	var testDirName = "test"
	var absTestDirName, _ = filepath.Abs(testDirName)
	var exists = directoryExist(absTestDirName)

	if exists {
		t.Error("Invalid result: directory should not exist")
	}
}

func TestDirectoryDoesntExistAgainstExistingFile(t *testing.T) {
	var testFilename = "test"
	var absTestFilename, _ = filepath.Abs(testFilename)
	var _, _ = os.Create(absTestFilename)
	defer os.Remove(absTestFilename)

	var exists = directoryExist(absTestFilename)

	if exists {
		t.Error("Invalid result: directory should not exist if there is a file with the same name")
	}
}

func TestFileExist(t *testing.T) {
	var testFilename = "test"
	var absTestFilename, _ = filepath.Abs(testFilename)
	var _, _ = os.Create(absTestFilename)
	defer os.Remove(absTestFilename)

	var exists = fileExist(absTestFilename)

	if !exists {
		t.Error("Invalid result: file should exist")
	}
}

func TestFileDoesntExist(t *testing.T) {
	var testFilename = "test"
	var absTestFilename, _ = filepath.Abs(testFilename)
	var exists = fileExist(absTestFilename)

	if exists {
		t.Error("Invalid result: file should not exist")
	}
}

func TestFileDoesntExistAgainstExistingDirectory(t *testing.T) {
	var testDirName = "test"
	var absTestDirName, _ = filepath.Abs(testDirName)
	var _ = os.Mkdir(absTestDirName, os.ModePerm)
	defer os.Remove(absTestDirName)

	var exists = fileExist(testDirName)

	if exists {
		t.Error("Invalid result: file should not exist if there is a directory with the same name")
	}
}
