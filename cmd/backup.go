// Copyright © 2016 Samsung CNCT
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Backup (ns, pod) to state store s (s3 bucket for now)
func Backup(ns string, pod string, s string) (err error) {
	return nil
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:          "backup [-s bucket]",
	Short:        "Backs up a GitLab deployment and saves the state to an s3 bucket",
	SilenceUsage: true,
	Long:         `Backs up `,
	Run: func(cmd *cobra.Command, args []string) {
		err := Backup("", "", "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			ExitCode = 1
			return
		}

		ExitCode = 0
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
}
