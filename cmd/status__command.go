package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [PATH]",
	Short: "Show status",
	Long:  `Show status of each git repository found in PATH`,

	Example: "gitas status ~ --name=p -b=true -o=n\ngitas status -lus\ngitas status /home --time=false",
	Aliases: []string{"ll"},

	Args: cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		statusMain(args)
	},
}

var config tConfig // Holds status' configuration

/*
init sets the flags
*/
func init() {
	rootCmd.AddCommand(statusCmd)

	initChoices()

	/* Init flags */

	statusCmd.Flags().SortFlags = false
	statusCmd.Flags().VarP(config.nameShown, "name", "n", "name shown: unique|path|short") // Choice

	statusCmd.Flags().BoolVarP(&config.showCommitTime, "time", "t", true, "time of last commit shown")
	statusCmd.Flags().VarP(config.timeFormat, "format", "f", "format time: relative|iso") // Choice

	statusCmd.Flags().BoolVarP(&config.showBranchHead, "branch", "b", false, "branch shown")
	statusCmd.Flags().BoolVarP(&config.showFetchNeeded, "query", "q", false, "query fetch needed (implies -br)")
	statusCmd.Flags().BoolVarP(&config.showBranchUpstream, "remote", "r", false, "remote shown")

	statusCmd.Flags().BoolVarP(&config.showUrl, "url", "l", false, "url shown")

	statusCmd.Flags().BoolVarP(&config.showDirty, "dirty", "d", true, "dirty shown")
	statusCmd.Flags().BoolVarP(&config.showUntracked, "untracked", "u", false, "untracked shown")
	statusCmd.Flags().BoolVarP(&config.showStash, "stash", "s", false, "stash shown")

	statusCmd.Flags().VarP(config.sortOrder, "order", "o", "order: time|name")                 // Choice
	statusCmd.Flags().VarP(config.emitFormat, "emit", "e", "emit format: table|json|markdown") // Choice
}

/*
Main status function

	'args' given command line arguments, that contain the root os search path
*/
func statusMain(args []string) {

	var givenDir string

	checkLogginglevel(args)

	/* Default the PATH */

	if len(args) != 1 {
		if loggingLevel >= 1 {
			logWarning.Println("Defaulting PATH")
		}
		givenDir = "."
	} else {
		givenDir = args[0]
	}

	/* Show branch infos when querying sync need */

	if config.showFetchNeeded {
		config.showBranchHead = true
		config.showBranchUpstream = true
	}

	/* Query all data when emitting json */

	if config.emitFormat.Value == "j" {
		config.showUrl = true
		config.showCommitTime = true
		config.showBranchHead = true
		config.showBranchUpstream = true
		config.showDirty = true
		config.showUntracked = true
		config.showStash = true
		config.timeFormat.Value = "I"
		config.showFetchNeeded = true
		config.showFetchNeeded = true
	}

	/* Get repos under 'givenDir' */

	repos, err := getReposDictionary(givenDir, config)
	if err != nil {
		logError.Fatalln(fmt.Errorf("getting repos dictionary failed. %w", err))
	}
	if loggingLevel >= 3 {
		logInfo.Printf("repos: %+v", repos)
	}

	/* Sort repositories */

	if config.sortOrder.Value == "n" {
		sort.SliceStable(repos, func(i, j int) bool {
			return repos[i].UniqueName < repos[j].UniqueName
		})
	} else {
		sort.SliceStable(repos, func(i, j int) bool {
			return repos[i].LastCommitEpoch > repos[j].LastCommitEpoch
		})
	}
	if loggingLevel >= 1 {
		logInfo.Println("repos sorted.")
	}

	/* Emit results */

	switch thisFormat := config.emitFormat.Value; thisFormat {
	case "j":
		if err := emitJson(repos); err != nil {
			logError.Fatalln(fmt.Errorf("emitting json failed. %w", err))
		}
	case "t":
		if err := emitTable(repos); err != nil {
			logError.Fatalln(fmt.Errorf("emitting table failed. %w", err))
		}
	case "m":
		emitMarkdown(repos)
	}

	if loggingLevel >= 1 {
		logInfo.Println("result emitted.")
	}

}
