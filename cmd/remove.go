/*
Copyright © 2020 Dimitar Ralev

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/dmralev/scaff/scaff"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [directory|filepath] [namespace]",
	Short: "Remove a filepath or namespace.",
	Long: `Remove given directory or filepath from a namespace, or the namespace itself.
If the command receives one argument, it is assumed that it has received namespace to delete.
Example:
  scaff remove test_namespace

Two provided arguments means that a [file or directory] should be looked up in the given [namespace] and removed from there.
Example
  scaff remove README.md test_namespace

While the files that scaff manages are supposed to be just a copy. It's useful to double check if the namespace files are important and need to be backed up in some way.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("Add requires a directory/file or a namespace argument.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var fileOrDir, namespace string
		if len(args) == 1 {
			fileOrDir, namespace = "", args[0]
		} else if len(args) == 2 {
			fileOrDir, namespace = args[0], args[1]
		}

		result, err := scaff.Remove(fileOrDir, namespace)
		if err != nil {
			result = err.Error()
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
