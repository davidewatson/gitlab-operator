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
	"os/exec"
	"strings"
	"time"
)

const maxApplyTimeout = 10 // Seconds

// Run command with args and kill if timeout is reached
func RunCommand(name string, args []string, timeout time.Duration) error {
	fmt.Printf("Running command \"%v %v\"\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)

	err := cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("Failed to kill command %v, err %v", name, err))
		}
		err = fmt.Errorf("Command %v timed out\n", name)
		break
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command %v returned err %v\n", name, err)
		}
		break
	}
	if err != nil {
		output, e := cmd.CombinedOutput()
		if e != nil {
			return e
		}
		fmt.Fprintf(os.Stderr, "%v", output)
		return err
	}
	fmt.Printf("Command %v completed successfully\n", name)

	return nil
}
