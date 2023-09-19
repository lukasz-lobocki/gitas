package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

/*
getRemoteSyncNeed returns if local needs sync with remote

	'thisRepo' structure (passed by refereferce) that contains initial data and to be populated
*/
func getRemoteSyncNeed(thisRepo *tRepo) error {

	commCommand := "git"
	repoUpstream := getStringRegex(`(.*)\/`, thisRepo.BranchUpstream) //name before '/'
	argsRemote := append([]string{},
		"remote", "show", repoUpstream,
	)

	cmdOutput, err := runCommand(commCommand, argsRemote, thisRepo.TopLevelPath)
	if err != nil {
		return fmt.Errorf("getting remote output failed. %w", err)
	}

	/* Build the output */

	if loggingLevel >= 3 {
		logInfo.Printf("git status command output for %v:\n%s\n", thisRepo, cmdOutput)
	}

	/* Construct boolean information */

	thisRepo.RemoteSyncNeed = !regexp.MustCompile(
		`(?mU)^\s*` +
			thisRepo.BranchHead +
			`\s*pushes to [[:print:]]*\s*\(` + UP_TO_DATE + `\)$`).
		MatchString(cmdOutput)
	return nil
}

/*
getOriginUrl returns remote's url

	'thisRepo' structure (passed by refereferce) that contains initial data and to be populated
*/
func getOriginUrl(thisRepo *tRepo) error {

	commCommand := "git"
	argsRemote := append([]string{},
		"config", "--get", "remote.origin.url",
	)

	cmdOutput, err := runCommand(commCommand, argsRemote, thisRepo.TopLevelPath)
	if err != nil {
		return fmt.Errorf("getting remote url failed. %w", err)
	}

	/* Build the output */

	if loggingLevel >= 3 {
		logInfo.Printf("git status command output for %v:\n%s\n", thisRepo, cmdOutput)
	}

	thisRepo.OriginUrl = cmdOutput

	return nil
}

/*
GetLastCommitTime returns the last commit time for the repo

	'thisRepo' structure (passed by refereferce) that contains initial data and to be populated
	'getEpoch' if result should be expressed in UNIX time
	'config' dictates pretty format for human readable time
*/
func getLastCommitTime(thisRepo *tRepo, getEpoch bool, config tConfig) error {

	var argsRemote []string

	commCommand := "git"
	if getEpoch {
		argsRemote = append([]string{},
			"log", "--pretty=format:%ct", "--date-order", "-n 1",
		)
	} else {
		argsRemote = append([]string{},
			"log", "--pretty=format:%c"+config.timeFormat.Value, "--date-order", "-n 1",
		)
	}

	cmdOutput, err := runCommand(commCommand, argsRemote, thisRepo.TopLevelPath)
	if err != nil {
		return fmt.Errorf("getting last commit time failed. %w", err)
	}

	/* Build the output */

	if loggingLevel >= 3 {
		logInfo.Printf("git status command output for %v:\n%s\n", thisRepo, cmdOutput)
	}

	if getEpoch {
		thisRepo.LastCommitEpoch = cmdOutput
	} else {
		thisRepo.LastCommitTime = cmdOutput
	}

	return nil
}

/*
getRepoStatus returns repo status

	'thisRepo' structure (passed by refereferce) that contains initial data and to be populated
	'config' dictates how data should be populated
*/
func getRepoStatus(thisRepo *tRepo, config tConfig) error {

	commCommand := "git"
	argsStatus := getArgsStatus(config)

	cmdOutput, err := runCommand(commCommand, argsStatus, thisRepo.TopLevelPath)
	if err != nil {
		return fmt.Errorf("getting status output failed. %w", err)
	}

	/* Build the output */

	if loggingLevel >= 3 {
		logInfo.Printf("git status command output for %v:\n%s\n", thisRepo, cmdOutput)
	}

	/* Construct ahead / behind symbol */

	regex := regexp.MustCompile(`(?mU)^# branch.ab \+(\d*) \-(\d*)$`) // Capture 2 numbers
	matches := regex.FindStringSubmatch(cmdOutput)
	if len(matches) > 0 {

		/* Convert both to int */

		if thisRepo.Ahead, err = strconv.Atoi(matches[1]); err != nil {
			return fmt.Errorf("converting ahead to int failed. %w", err)
		}
		if thisRepo.Behind, err = strconv.Atoi(matches[2]); err != nil {
			return fmt.Errorf("converting behind to int failed. %w", err)
		}

		/* Calculate StatusAB using ahead & behind values */

		if thisRepo.Ahead == 0 && thisRepo.Behind == 0 {
			thisRepo.StatusAB = SYNCED_CHAR
		} else if thisRepo.Ahead == 0 {
			thisRepo.StatusAB = REMOTE_AHEAD_CHAR
		} else if thisRepo.Behind == 0 {
			thisRepo.StatusAB = LOCAL_AHEAD_CHAR
		} else {
			thisRepo.StatusAB = DIVERGED_CHAR
		}
	}

	/* Construct boolean information */

	if config.showDirty {
		thisRepo.Dirty = regexp.MustCompile(`(?m)^[^#?]`).MatchString(cmdOutput) // Dirty, does not beging with `#` nor `?`
	}
	if config.showUntracked {
		thisRepo.Untracked = regexp.MustCompile(`(?m)^\?`).MatchString(cmdOutput) // Untracked, begins with `?`
	}
	if config.showStash {
		thisRepo.Stash = regexp.MustCompile(`(?mU)^# stash (\d*)$`).MatchString(cmdOutput)
	}

	/* Construct branch information */

	if config.showBranchHead || config.showRemoteSyncNeed {
		thisRepo.BranchHead = getStringRegex(`(?mU)^# branch.head (.+)$`, cmdOutput)
	}
	if config.showBranchUpstream || config.showUrl || config.showRemoteSyncNeed {
		thisRepo.BranchUpstream = getStringRegex(`(?mU)^# branch.upstream (.+)$`, cmdOutput)
	}

	return nil

}

func runCommand(thisCommand string, thisArgs []string, thisDir string) (string, error) {
	cmd := exec.Command(thisCommand, thisArgs...)
	cmd.Dir = thisDir
	out_bytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting cmd output failed. %w", err)
	}

	return strings.TrimSpace(string(out_bytes)), nil
}

/*
getArgsStatus returns git status command arguments

	'config' rules of the construction

Convenience funtion
*/
func getArgsStatus(config tConfig) []string {
	args := append([]string{},
		"status", "--branch", "--porcelain=2",
	)
	if config.showStash {
		args = append(args,
			"--show-stash",
		)
	}
	if !config.showUntracked {
		args = append(args,
			"--untracked-files=no",
		)
	}
	if config.lookForSubGits {
		args = append(args,
			"--ignore-submodules=all",
		)
	}
	return args
}

/*
getStringRegex returns first regex 'expression' match within the 'input'

	'expression' regex expression
	'input' string to be searched
*/
func getStringRegex(expression string, input string) string {
	regex := regexp.MustCompile(expression)
	matches := regex.FindStringSubmatch(input)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}
