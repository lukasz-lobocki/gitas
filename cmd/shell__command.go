package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell [PATH] \"command\"",
	Short: "Execute command",
	Long:  `Execute "command" for each git repository found in PATH`,

	Example: "gitas shell /home \"ls\"\ngitas shell ~ \"git describe --abbrev=0 --tags\"\ngitas shell \"ls | grep 'P'\"",

	Args: cobra.RangeArgs(1, 2),

	Run: func(cmd *cobra.Command, args []string) {
		shellMain(args)
	},
}

// Cobra initiation
func init() {
	rootCmd.AddCommand(shellCmd)
}

/*
Shell main function

	'args' given command line arguments, that contain the command to be run by shell
*/
func shellMain(args []string) {
	var (
		cmdArgs  []string // Args of thecommand to execute
		givenDir string
		err      error
	)

	checkLogginglevel(args)

	/* Construct arguments */

	switch lenArgs := len(args); lenArgs {
	case 1:
		givenDir = "."
		cmdArgs = append([]string{"-c"}, args[0])
	case 2:
		givenDir = args[0]
		cmdArgs = append([]string{"-c"}, args[1])
	}

	/* Find repos */

	git_slice, err := findRepos(givenDir, config.lookForSubGits)
	if err != nil {
		logError.Fatalln(fmt.Errorf("finding repos failed. %w", err))
	}
	if loggingLevel >= 2 {
		logInfo.Printf("%d repos found.", len(git_slice))
	}

	/* Execute for each repo */

	for _, thisGit := range git_slice {
		execShell(thisGit, cmdArgs)
	}

}

/*
execShell spawns shell to run arbitrary command within given path

	'thisGit' path
	'shellCommand' command passed to shell
*/
func execShell(thisGit string, shellCommand []string) {

	cmd := exec.Command(SHELL, shellCommand...)

	/* 	Pipe the commands output to the applications standard output */

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = thisGit

	if loggingLevel >= 3 {
		logInfo.Printf("execShell: in %s starting %s %s \"%s\"", thisGit, SHELL, shellCommand[0], shellCommand[1])
	}

	/* Actual run */

	if err := cmd.Run(); err != nil {
		if loggingLevel >= 1 {
			logWarning.Printf("error %s, running %s %s \"%s\" in %s\n", err, SHELL, shellCommand[0], shellCommand[1], thisGit)
		}
	}

	if loggingLevel >= 3 {
		logInfo.Printf("shell finished")
	}
}
