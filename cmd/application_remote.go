package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var remoteName string

var appRemoteCmd = &cobra.Command{
	Use:     "remote",
	Short:   "Add a remote to your current git repository",
	Example: "civo app remote APP_NAME -r REMOTE_NAME",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		utility.EnsureCurrentRegion()
		client, err := config.CivoAPIClient()

		if regionSet != "" {
			client.Region = regionSet
		}

		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("App %s not found", findApp.Name)
			os.Exit(1)
		}

		// Ensure current folder is a git repository.
		status_cmd := exec.Command("git", "status")
		status_cmd.Stdout = &stdout
		status_cmd.Stderr = &stderr
		err = status_cmd.Run()
		if err != nil {
			fmt.Printf("error running: git %s\n", strings.Join(status_cmd.Args, " "))
			fmt.Println(stderr.String())
		} else {
			if stdout.String() != "fatal: not a git repository (or any of the parent directories): .git\n" {
				fmt.Println("")
			} else {
				utility.Error("You must be in a git repository to use this command")
				os.Exit(1)
			}
		}

		// Check if current remote exists.
		remote_cmd := exec.Command("git", "ls-remote")
		remote_cmd.Stdout = &stdout
		remote_cmd.Stderr = &stderr
		err = remote_cmd.Run()
		if err != nil {
			fmt.Printf("error running: git %s\n", strings.Join(remote_cmd.Args, " "))
			fmt.Println(stderr.String())
		}

		//Add a new remote.
		remote_cmd = exec.Command("git", "remote", "add", remoteName, fmt.Sprintf("git@git.civo.app:%s/%s/%s", client.Region, client.GetAccountID(), findApp.Name))
		//How to get account id from api key above?
		remote_cmd.Stdout = &stdout
		remote_cmd.Stderr = &stderr
		err = remote_cmd.Run()
		if err != nil {
			fmt.Printf("error running: git %s\n", strings.Join(remote_cmd.Args, " "))
			fmt.Println(stderr.String())
		} else {
			fmt.Printf("Added the \"%s\" remote to your git repository", remoteName)
		}
	},
}
