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
	"os"
	"time"

	"github.com/spf13/cobra"
)

// Find the one expected pod with the label selector in this namespace and run
// GitLab backup on it. Store the result in an s3 bucket.
func Backup(s3 string) error {
	namespace, err := GetNamespace()
	if err != nil {
		return err
	}

	key, value := GitLabLabelKey, GitLabLabelValue
	podNames, err := GetPodsWithLabel(namespace, key, value)
	if err != nil {
		return err
	} else if len(podNames) != 1 {
		return fmt.Errorf("there were %v pods with label %v=%v when exactly one was expected\n", len(podNames), key, value)
	}

	fmt.Printf("Begining backup of GitLab instance %v, %v\n", namespace, podNames[0])

	options := ExecOptions{
		Command:       nil,
		Namespace:     namespace,
		PodName:       podNames[0],
		ContainerName: GitLabContainerName,
		CaptureStdout: true,
		CaptureStderr: true,
	}

	options.Command = []string{"gitlab-rake", "gitlab:backup:create"}
	err = ExecWithOptions(options)
	if err != nil {
		return err
	}

	filename := GitLabBackupPrefix + time.Now().UTC().Format(time.RFC3339) + ".tar.gz"

	options.Command = []string{"tar", "czf", filename, "/etc/gitlab"}
	err = ExecWithOptions(options)
	if err != nil {
		return err
	}

	options.Command = []string{"aws", "s3", "cp", filename, s3}
	err = ExecWithOptions(options)
	if err != nil {
		return err
	}

	fmt.Printf("Finished backup of GitLab instance\n")

	return nil
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:          "backup [-s bucket]",
	Short:        "Backs up GitLab",
	SilenceUsage: true,
	Long:         `Backs up a GitLab deployment and saves the state to an s3 bucket.`,
	Run: func(cmd *cobra.Command, args []string) {
		s3 := operatorConfig.GetString("s3")
		err := Backup(s3)
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
