/*
Copyright © 2019 Civo Ltd <hello@civo.com>

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
package main

import (
	"fmt"
	"os"

	"github.com/civo/cli/cmd"
	"github.com/civo/cli/common"
	"github.com/savioxavier/termlink"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			resp, skip := common.VersionCheck()
			gitIssueLink := termlink.ColorLink("GitHub issue", "https://github.com/civo/cli/issues", "italic green")
			if skip == true {
				fmt.Printf("panic : %s \nPlease check if you are using the latest version of CLI and retry the command \nIf you are still facing issues, please report it on our community slack or open a %s \n", err, gitIssueLink)
				os.Exit(1)
			}
			res := resp.Current
			updateCmd := "civo update"
			fmt.Printf("panic : %s \nYour CLI Version : %s \nPlease, run %q and retry the command \nIf you are still facing issues, please report it on our community slack or open a %s \n", err, res, updateCmd, gitIssueLink)
		}
	}()
	cmd.Execute()
}
