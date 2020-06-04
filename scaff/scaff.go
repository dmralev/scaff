package scaff

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

var filePaths []string

var home, homeErr = homedir.Dir()
var namespaceHome = path.Join(home, ".scaff", "namespaces")

// Initialize base settings
// Create the needed directories for storing namespaces
func Init() error {
	if homeErr != nil {
		return homeErr
	}

	_, missingOrOther := os.Stat(namespaceHome)
	if missingOrOther != nil {
		err := os.MkdirAll(namespaceHome, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

// Add files to a given namespace. Grab all filepaths from a given source
// and write their contents into a new files with the same names and folder structure
// under a common namespace.
func Add(src, namespace string) (string, error) {
	var fileCount int
	var dirCount int

	srcInfo, err := os.Stat(src)
	if err != nil {
		errMessage := fmt.Sprintf("Error: %s no such file or directory.", src)
		return "", errors.New(errMessage)
	}

	if srcInfo.IsDir() {
		files, _ := ioutil.ReadDir(src)
		if len(files) < 1 {
			return "", errors.New("Error: Chosen directory is empty.")
		}
	}

	namespaceDir := path.Join(namespaceHome, namespace)
	os.Mkdir(namespaceDir, 0777)
	if err != nil {
		return "", err
	}

	err = filepath.Walk(src, func(pathname string, info os.FileInfo, err error) error {
		filename := path.Base(pathname)
		relPath := strings.TrimPrefix(pathname, src)

		if strings.HasPrefix(filename, ".") || strings.HasPrefix(relPath, "/.") {
			return nil
		}

		if pathname == src {
			_, relPath = path.Split(pathname)
		}

		if info.IsDir() {
			// Exclude the root folder
			if src == pathname {
				return nil
			}

			// TODO: Copy the permissions
			newDir := path.Join(namespaceDir, relPath)
			err := os.Mkdir(newDir, 0777)
			if err != nil {
				return err
			}
			dirCount += 1
			return nil
		}

		dest := path.Join(namespaceDir, relPath)

		contents, err := ioutil.ReadFile(pathname)
		if err != nil {
			return err
		}

		// TODO: Copy the permissions
		err = ioutil.WriteFile(dest, contents, 0777)
		if err != nil {
			return err
		}
		fileCount += 1

		return nil
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d files and %d directories stored under namespace %s", fileCount, dirCount, namespace), nil
}

// Remove a whole namespace or a single file/folder.
func Remove(delPath, namespace string) (string, error) {
	namespaceDir := path.Join(namespaceHome, namespace)
	delPath = path.Join(namespaceDir, delPath)

	_, err := os.Stat(namespaceDir)
	if err != nil {
		errMessage := fmt.Sprintf("Error: Namespace %s not found", namespace)
		return "", errors.New(errMessage)
	}

	_, err = os.Stat(delPath)
	if err != nil {
		errMessage := fmt.Sprintf("Error: Path not found in %s namespace", namespace)
		return "", errors.New(errMessage)
	}

	// TODO: Extract into func
	fileCount, dirCount := 0, 0
	filepath.Walk(delPath, func(pathname string, info os.FileInfo, err error) error {
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
	if namespaceDir == delPath {
		fmt.Printf("Removing namespace %s, are you sure? [y/n]\n", namespace)
	} else {
		fmt.Printf("About to remove %d files and %d directories from namespace %s, are you sure? [y/n]\n", fileCount, dirCount, namespace)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	switch input {
	case "y":
		err = os.RemoveAll(delPath)
		if err != nil {
			return "", err
		}

		if namespaceDir == delPath {
			return fmt.Sprintf("Removed namespace %s.", namespace), nil
		}

		return fmt.Sprintf("Removed %d files and %d directories from namespace %s.", fileCount, dirCount, namespace), nil
	case "n":
		return "Remove cancelled.", nil
	default:
		fmt.Println(input)
		return "Unexpected input, remove cancelled.", nil
	}
}

// List all namespaces along with a short stats
func List() (string, error) {
	dirs, err := ioutil.ReadDir(namespaceHome)
	if err != nil {
		return "", err
	}

	// Init writer
	// params after stdout - minwidth, tabwidth, padding, padchar, flags
	b := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(b, 5, 4, 4, ' ', 0)
	rowInfo := "%s\t%d folders\t%d files"

	// Setup the first rows
	fmt.Println()
	fmt.Fprintln(w, "Namespaces\tStats\t")
	fmt.Fprintln(w, "---\t---\t")

	// Count and show stats
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}
		fileCount, dirCount := 0, 0
		namespaceDir := path.Join(namespaceHome, dir.Name())
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

	return b.String(), nil
}

// Show tree structure of the files in a given namespace
func Show(namespace string) (string, error) {
	namespaceDir := path.Join(namespaceHome, namespace)

	tree, err := Tree(namespaceDir, "")
	if err != nil {
		return "", err
	}

	return tree, nil
}

// Copy the files from a namespace, to a given directory
func Get(dest, namespace string) (string, error) {
	_, err := os.Stat(dest)
	if err != nil {
		return "", err
	}

	// Returns the path without the last part <3
	namespaceDir := path.Join(namespaceHome, namespace)
	_, err = os.Stat(namespaceDir)
	if err != nil {
		errMessage := fmt.Sprintf("Error: Namespace %s not found", namespace)
		return "", errors.New(errMessage)
	}

	dirCount, fileCount := 0, 0
	filepath.Walk(namespaceDir, func(pathname string, info os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(pathname, namespaceDir)

		if info.IsDir() {
			// Try to exclude the root folder
			if namespaceDir == pathname {
				return nil
			}

			newDir := path.Join(dest, relPath)
			// TODO: Copy the permissions
			os.Mkdir(newDir, 0777)
			dirCount += 1

			return nil
		}

		destDir := path.Join(dest, relPath)

		contents, err := ioutil.ReadFile(pathname)
		if err != nil {
			return err
		}
		fileCount += 1

		// TODO: Copy the permissions
		err = ioutil.WriteFile(destDir, contents, 0777)
		if err != nil {
			return err
		}

		return nil
	})

	return fmt.Sprintf("%d files and %d directories copied from namespace %s", fileCount, dirCount, namespace), nil
}

func Tree(namespaceDir, prefix string) (string, error) {
	buffer := bytes.NewBufferString("")
	nodes, err := ioutil.ReadDir(namespaceDir)
	if err != nil {
		return "", err
	}

	for index, node := range nodes {
		// Print symbols depending on if the current file is last in the directory or not
		if index == len(nodes)-1 {
			fmt.Fprintln(buffer, prefix+"└──"+" "+node.Name())
		} else {
			fmt.Fprintln(buffer, prefix+"├──"+" "+node.Name())
		}

		// Go deeper if the node
		// Print symbol if there are more directories or files following it
		if node.IsDir() {
			nodeDir := path.Join(namespaceDir, node.Name())
			if index < len(nodes)-1 {
				Tree(nodeDir, prefix+"│   ")
			} else {
				Tree(nodeDir, prefix+"    ")
			}
		}
	}

	return buffer.String(), nil
}
