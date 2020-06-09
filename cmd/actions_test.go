package cmd

import (
	"bytes"
	"fmt"
	"github.com/dmralev/scaff/scaff"
	homedir "github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var wd, _ = os.Getwd()
var parentDir = path.Dir(wd)

var home, homeErr = homedir.Dir()
var namespaceHome = path.Join(home, ".scaff", "namespaces")

func AddTestFile() {
	scaff.Init()

	// TODO: Needs nested directories
	testFile := path.Join(parentDir, "LICENSE")
	scaff.Add(testFile, "test")
}

func assertEqualStrings(x, y string, t *testing.T) {
	if x == y {
		return
	}
	t.Errorf("Expected %s but got %s", x, y)
}

func assertContains(x, y string, t *testing.T) {
	if strings.Contains(x, y) {
		return
	}
	t.Errorf("%s doesn't contain %s", x, y)
}

func clearTest() {
	testNamespace := path.Join(namespaceHome, "test")
	os.RemoveAll(testNamespace)
}

//
func TestAddFile(t *testing.T) {
	filePath := path.Join(parentDir, "LICENSE")
	buffer := bytes.NewBufferString("")

	rootCmd.SetArgs([]string{"add", filePath, "test"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	t.Cleanup(clearTest)
}

// Basic test in which we don't care about
// the false positives that the assertion can give
func TestAddDir(t *testing.T) {
	dirPath := path.Join(parentDir, "cmd")
	buffer := bytes.NewBufferString("")

	rootCmd.SetArgs([]string{"add", dirPath, "test"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	testNamespace := path.Join(namespaceHome, "test")

	// Assert that namespace exists
	srcFiles, _ := ioutil.ReadDir(dirPath)
	destFiles, err := ioutil.ReadDir(testNamespace)
	if err != nil {
		t.FailNow()
	}

	// Let this snippet be remembered as the moment I found out
	// that Go doesn't have .Contains for slices
	// TODO: Should I make a helper file and extract such logic there?
	fileCount, dirCount := 0, 0
	for _, srcFile := range srcFiles {
		isFound := false
		if strings.HasPrefix(srcFile.Name(), ".") || strings.HasPrefix(srcFile.Name(), "/.") {
			continue
		}
		for _, destFile := range destFiles {
			if srcFile.Name() == destFile.Name() {
				isFound = true
				if srcFile.IsDir() {
					dirCount += 1
				} else {
					fileCount += 1
				}
			}
		}
		if !isFound {
			t.FailNow()
		}
	}

	expectedString := fmt.Sprintf("%d files and %d directories stored under namespace %s\n", fileCount, dirCount, "test")
	assertEqualStrings(buffer.String(), expectedString, t)

	t.Cleanup(clearTest)
}

//
func TestAddValidation(t *testing.T) {
	errMessage := "Error: Add requires a directory/file and a namespace arguments."
	buffer := bytes.NewBufferString("")
	dirPath := path.Join(parentDir, "cmd")

	// Test with single param
	rootCmd.SetArgs([]string{"add", dirPath})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	assertContains(buffer.String(), errMessage, t)

	// Test with more than two params
	buffer = bytes.NewBufferString("")

	rootCmd.SetArgs([]string{"add", dirPath, "test", "mistake", "another", "other"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	assertContains(buffer.String(), errMessage, t)

	t.Cleanup(clearTest)
}

//
func TestRemoveValidation(t *testing.T) {
	errMessage := "Error: Add requires a directory/file or a namespace argument."
	buffer := bytes.NewBufferString("")

	AddTestFile()
	dirPath := path.Join(parentDir, "cmd")

	// Test with single param
	rootCmd.SetArgs([]string{"remove"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	assertContains(buffer.String(), errMessage, t)

	// Test with more than required params
	buffer = bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"remove", dirPath, "test", "mistake", "another", "other"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	assertContains(buffer.String(), errMessage, t)

	t.Cleanup(clearTest)
}

// TODO Needs research on how to write to stdin when prompted
func TestRemove(t *testing.T) {
	// AddTestFile()
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
	// clearTest()
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
	AddTestFile()
	buffer := bytes.NewBufferString("")

	rootCmd.SetArgs([]string{"list"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	assertContains(buffer.String(), "test", t)

	t.Cleanup(clearTest)
}

// Basic test
func TestShow(t *testing.T) {
	// TODO: Needs case for nested directories
	AddTestFile()
	buffer := bytes.NewBufferString("")

	rootCmd.SetArgs([]string{"show", "test"})
	rootCmd.SetOut(buffer)
	rootCmd.Execute()

	testNamespace := path.Join(namespaceHome, "test")

	// Test if all files can be found in the tree string
	// TODO: It can give false positive when there is more than one file with the same name
	namespaceFiles, _ := ioutil.ReadDir(testNamespace)
	for _, file := range namespaceFiles {
		assertContains(buffer.String(), file.Name(), t)
	}

	// Test actual pretty basic Tree formatting by directly calling Tree, no other way for now
	expectedList, _ := scaff.Tree(testNamespace, "")
	assertEqualStrings(buffer.String(), expectedList, t)

	t.Cleanup(clearTest)
}

// func TestGet(t *testing.T) {
// 	Get("/Users/dimitarralev/code/testee", "dimitarralev")
// }
