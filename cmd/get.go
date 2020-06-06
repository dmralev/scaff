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
	"os"

	"github.com/dmralev/scaff/scaff"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [namespace]",
	Short: "Copy files from a given namespace to the current directory",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Get requires a namespace argument only.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		namespace := args[0]

		result, err := scaff.Get(wd, namespace)
		if err != nil {
			result = err.Error()
		}

		fmt.Fprintln(cmd.OutOrStdout(), result)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
