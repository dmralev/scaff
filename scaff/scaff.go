package scaff

import (
	// "fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var workDir, err = os.Getwd()
var currentDir = path.Dir(workDir)
var namespaceHome = "namespaces"

var filePaths []string

// Add files to a given namespace. Grab all filepaths from a given source
// and write their contents into a new files with the same names and folder structure
// under a common namespace.
func Add(src, namespace string) bool {

	_, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// Returns the path without the last part <3
	namespaceDir := path.Join(currentDir, namespaceHome, namespace)
	os.Mkdir(namespaceDir, 0777)

	filepath.Walk(src, func(pathname string, info os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(pathname, src)

		if info.IsDir() {
			// Try to exclude the root folder
			if src == pathname {
				return nil
			}

			newDir := path.Join(namespaceDir, relPath)
			os.Mkdir(newDir, 0777)
			return nil
		}

		destDir := path.Join(namespaceDir, relPath)

		contents, err := ioutil.ReadFile(pathname)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(destDir, contents, 0777)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	return true
}

// Remove a whole namespace or a single file/folder.
func Remove(relPath, namespace string) bool {
	namespaceDir := path.Join(currentDir, namespaceHome, namespace)

	// Check if exists
	// TODO: You definitely want to (friendly)log to the user if there is no such folder
	_, err := os.Stat(namespaceDir)
	if err != nil {
		log.Fatal(err)
		return false
	}

	delPath := path.Join(namespaceDir, relPath)
	err = os.RemoveAll(delPath)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// TODO Log what was succesfully deleted

	return true
}

func List(namespace string) string {
	return "Not implemented"
}

func Show(namespace string) string {
	return "Not implemented"
}

func Use(namespace string) string {
	// TODO or Get?
	return "Not implemented"
}
