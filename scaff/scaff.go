package scaff

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
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

func List() {
	namespaceRoot := path.Join(currentDir, namespaceHome)

	dirs, err := ioutil.ReadDir(namespaceRoot)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Init writer
	// params after stdout - minwidth, tabwidth, padding, padchar, flags
	w := tabwriter.NewWriter(os.Stdout, 5, 4, 4, ' ', 0)
	rowInfo := "%s\t%d folders\t%d files"

	// Setup the first rows
	fmt.Println()
	fmt.Fprintln(w, "Namespaces\tStats\t")
	fmt.Fprintln(w, "---\t---\t")

	// Count and show stats
	for _, dir := range dirs {
		fileCount, dirCount := 0, 0
		namespaceDir := path.Join(namespaceRoot, dir.Name())
		filepath.Walk(namespaceDir, func(pathname string, info os.FileInfo, err error) error {
			// Don't count namespace itself as folder
			if pathname == namespaceDir {
				return nil
			}

			if info.IsDir() {
				dirCount += 1
			} else {
				fileCount += 1
			}

			return nil
		})

		row := fmt.Sprintf(rowInfo, dir.Name(), dirCount, fileCount)
		fmt.Fprintln(w, row)
	}
	w.Flush()

	return
}

func Show(namespace string) string {
	return "Not implemented"
}

func Use(namespace string) string {
	// TODO or Get?
	return "Not implemented"
}
