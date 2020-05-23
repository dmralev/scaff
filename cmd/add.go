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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [directory|filepath] [namespace]",
	Short: "Store a file or directory under a single namespace",
	Long:  `Add expects to receive a path to directory or file, and namespace under which to store the added files.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Add requires a directory/file and a namespace arguments.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fileOrDir, namespace := args[0], args[1]

		result, err := scaff.Add(fileOrDir, namespace)
		if err != nil {
			result = err.Error()
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
