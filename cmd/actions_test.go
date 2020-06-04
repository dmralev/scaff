package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/dmralev/scaff/scaff"
	homedir "github.com/mitchellh/go-homedir"
)

// Quickfix
var wd, _ = os.Getwd()
var parentDir = path.Dir(wd)

var home, homeErr = homedir.Dir()
var namespaceHome = path.Join(home, ".scaff", "namespaces")

func prepareAdd() {
	// Quickfix
	scaff.Init()

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

//
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

//
func TestRemoveValidation(t *testing.T) {
	prepareAdd()
	dirPath := path.Join(parentDir, "cmd")

	// Test with single param
	rootCmd.SetArgs([]string{"remove"})
	buffer := bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()

	errMessage := "Error: Add requires a directory/file or a namespace argument."
	if !strings.Contains(buffer.String(), errMessage) {
		t.Errorf(buffer.String())
	}

	rootCmd.SetArgs([]string{"remove", dirPath, "test", "mistake", "another", "other"})
	buffer = bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()
	if !strings.Contains(buffer.String(), errMessage) {
		t.Errorf(buffer.String())
	}

	clearAdd()

	// TODO: Assert files are removed
}

// TODO Needs research on how to write to stdin when prompted
func TestRemove(t *testing.T) {
	// prepareAdd()
	//
	// // Test with single file
	// rootCmd.SetArgs([]string{"remove", "LICENSE", "test"})
	// outBuffer := bytes.NewBufferString("")
	// inBuffer := bytes.NewBufferString("")
	// rootCmd.SetOut(outBuffer)
	// rootCmd.SetIn(inBuffer)
	//
	// inBuffer.WriteString("n")
	// rootCmd.Execute()
	// fmt.Println("exec - > ", outBuffer.String())
	// outBuffer.WriteString("n")
	// fmt.Println("after write - >", outBuffer.String())
	//
	// clearAdd()
	//
	// // Test with directory
	// // rootCmd.SetArgs([]string{"remove"})
	// // buffer := bytes.NewBufferString("")
	// // rootCmd.SetOut(buffer)
	// //
	// // rootCmd.Execute()
}

// Super basic test
func TestList(t *testing.T) {
	prepareAdd()

	rootCmd.SetArgs([]string{"list"})

	buffer := bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()

	if !strings.Contains(buffer.String(), "test") {
		t.FailNow()
	}

	clearAdd()
}

// Basic test
func TestShow(t *testing.T) {
	// TODO: Needs case for nested directories
	prepareAdd()

	rootCmd.SetArgs([]string{"show", "test"})

	buffer := bytes.NewBufferString("")
	rootCmd.SetOut(buffer)

	rootCmd.Execute()

	testNamespace := path.Join(namespaceHome, "test")

	// Test if all files can be found in the tree string
	// TODO: It can give false positive when there is more than one file with the same name
	namespaceFiles, _ := ioutil.ReadDir(testNamespace)
	for _, file := range namespaceFiles {
		if !strings.Contains(buffer.String(), file.Name()) {
			t.FailNow()
		}
	}

	// Test actual pretty basic Tree formatting by directly calling Tree, no other way for now
	expectedList, _ := scaff.Tree(testNamespace, "")
	if buffer.String() != expectedList {
		t.FailNow()
	}

	// TODO: Use the golang Clean method?
	clearAdd()
}

// func TestGet(t *testing.T) {
// 	Get("/Users/dimitarralev/code/testee", "dimitarralev")
// }
