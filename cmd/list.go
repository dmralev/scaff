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
	"fmt"
	"github.com/dmralev/scaff/scaff"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored namespaces",
	Long: `See a simple rundown with basic information about your namespaces.
Paths starting with . are excluded for now.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := scaff.List()
		if err != nil {
			result = err.Error()
		}

		fmt.Fprintf(cmd.OutOrStdout(), result)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
