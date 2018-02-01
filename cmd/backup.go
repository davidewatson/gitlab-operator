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
	"time"
)

func ConstructKubeExecCmdArgs(kubeConfig string, namespace string, podName string, containerName string) []string {
	return []string{"--kubeconfig", kubeConfig,
		"--namespace", namespace,
		"exec",
		podName,
		"-c",
		containerName,
		"--",
	}
}

// Backup (ns, pod) to state store s (s3 bucket for now)
func Backup(kubeConfig string, namespace string, podName string, containerName string, s3 string) (err error) {
	cmdName := "kubectl"
	cmdBaseArgs := ConstructKubeExecCmdArgs(kubeConfig, namespace, podName, containerName)
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	cmdExecArgs := []string{"echo", "gitlab-rake", "gitlab:backup:create"}
	cmdArgs := append(cmdBaseArgs, cmdExecArgs...)
	err = RunCommand(cmdName, cmdArgs, cmdTimeout)
	if err != nil {
		return err
	}

	cmdExecArgs = []string{"echo", "tar", "cfz", "etc-gitlab.tar.gz", "/etc/gitlab/"}
	cmdArgs = append(cmdBaseArgs, cmdExecArgs...)
	err = RunCommand(cmdName, cmdArgs, cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:          "backup [-s bucket]",
	Short:        "Backs up a GitLab deployment and saves the state to an s3 bucket",
	SilenceUsage: true,
	Long:         `Backs up `,
	Run: func(cmd *cobra.Command, args []string) {
		kubeConfig := operatorConfig.GetString("kubeconfig")
		namespace := operatorConfig.GetString("namespace")
		podName := operatorConfig.GetString("pod")
		containerName := operatorConfig.GetString("container")
		s3bucket := operatorConfig.GetString("s3")
		err := Backup(kubeConfig, namespace, podName, containerName, s3bucket)
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
