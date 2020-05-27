package cmd

import (
	"bytes"
	"github.com/dmralev/scaff/scaff"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var wd, _ = os.Getwd()
var parentDir = path.Dir(wd)
var namespaceHome = path.Join(parentDir, "namespaces")
var testNamespace = path.Join(namespaceHome, "test")

func prepareAdd() {
	// TODO: Needs nested directories
	testFile := path.Join(parentDir, "LICENSE")
	scaff.Add(testFile, "test")
}

func clearAdd() {
	testNamespace := path.Join(namespaceHome, "test")
	os.RemoveAll(testNamespace)
}

//
func TestAddFile(t *testing.T) {
	filePath := path.Join(parentDir, "LICENSE")

	rootCmd.SetArgs([]string{"add", filePath, "test"})

	buffer := bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()

	clearAdd()
}

//
func TestAddDir(t *testing.T) {
	dirPath := path.Join(parentDir, "cmd")

	rootCmd.SetArgs([]string{"add", dirPath, "test"})

	buffer := bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()

	// TODO: It can give false positive when there is more than one file with the same name
	// TODO: DRY the checks into a separate func?
	testNamespace := path.Join(namespaceHome, "test")

	srcFiles, _ := ioutil.ReadDir(dirPath)
	destFiles, err := ioutil.ReadDir(testNamespace)
	if err != nil {
		t.FailNow()
	}

	// Let this snippet be remembered as the moment I found out
	// that Go doesn't have .Contains for slices
	// TODO: Should I make a helper file and extract such logic there?
	for _, srcFile := range srcFiles {
		isFound := false
		for _, destFile := range destFiles {
			if srcFile.Name() == destFile.Name() {
				isFound = true
			}
		}
		if !isFound {
			t.FailNow()
		}
	}

	// TODO: Count the files and dirs above, to check the print message as well
	// You have the data you need

	clearAdd()
}

func TestAddValidation(t *testing.T) {
	dirPath := path.Join(parentDir, "cmd")

	// Test with single param
	rootCmd.SetArgs([]string{"add", dirPath})
	buffer := bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()

	errMessage := "Error: Add requires a directory/file and a namespace arguments."
	if !strings.Contains(buffer.String(), errMessage) {
		t.Errorf(buffer.String())
	}

	// Test with more than two params
	rootCmd.SetArgs([]string{"add", dirPath, "test", "mistake", "another", "other"})
	buffer = bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()
	if !strings.Contains(buffer.String(), errMessage) {
		t.Errorf(buffer.String())
	}

	clearAdd()
}
