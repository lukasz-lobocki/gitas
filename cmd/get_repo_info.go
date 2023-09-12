package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

/*
getOriginUrl returns remote's url

	'dirName' path to repo
*/
func getOriginUrl(dirName string) (string, error) {

	var (
		dir fs.FileInfo
		err error

		cmd       *exec.Cmd // Command to execute
		out_bytes []byte    //result of command execution
		out       string    // result of command execution
	)

	dir, err = os.Stat(dirName)
	if err != nil {
		return "", fmt.Errorf("getting stat failed. %w", err)
	}

	if dir.IsDir() {

		/* Execute the command */

		cmd = exec.Command("git", "config", "--get", "remote.origin.url")
		cmd.Dir = dirName
		out_bytes, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("getting cmd output failed. %w", err)
		}

		/* Build the output */

		out = strings.TrimSpace(string(out_bytes))

		if loggingLevel >= 3 {
			logInfo.Printf("\"%s\" git origin url output for %s\n", out, dirName)
		}

		return out, nil
	}

	return "", errors.New("not a directory")
}

/*
GetLastCommitTime returns the last commit time for the repo

	'dirName' path to the repo
	'getEpoch' if result should be expressed in UNIX time
	'config' dictates pretty format for human readable time
*/
func getLastCommitTime(dirName string, getEpoch bool, config tConfig) (string, error) {

	var (
		dir fs.FileInfo
		err error

		cmd       *exec.Cmd // Command to execute
		out_bytes []byte    // result of command execution
		out       string    // result of command execution
	)

	dir, err = os.Stat(dirName)
	if err != nil {
		return "", fmt.Errorf("getting stat failed. %w", err)
	}

	if dir.IsDir() {

		if getEpoch {
			cmd = exec.Command("git", "log", "--pretty=format:%ct", "--date-order", "-n 1")
		} else {
			cmd = exec.Command("git", "log", "--pretty=format:%c"+config.timeFormat.Value, "--date-order", "-n 1")
		}

		/* Execute the command */

		cmd.Dir = dirName
		out_bytes, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("getting cmd output failed. %w", err)
		}

		/* Build the output */

		out = strings.TrimSpace(string(out_bytes))

		if loggingLevel >= 3 {
			logInfo.Printf("\"%s\" git date command output for %s\n", out, dirName)
		}

		return out, nil
	}

	return "", errors.New("not a directory")
}

/*
getRepoStatus returns repo status

	'thisRepo' structure (passed by refereferce) that contains initial data and to be populated
	'config' dictates how data should be populated
*/
func getRepoStatus(thisRepo *tRepo, config tConfig) error {

	var (
		dir fs.FileInfo // Directory
		err error       // Errror

		cmd       *exec.Cmd // Command to execute
		out_bytes []byte    // Result of command execution
		cmdOutput string    // Result of command execution

		regex   *regexp.Regexp // Regular expression formula
		matches []string       // Result of regular expression matching
	)

	dir, err = os.Stat(thisRepo.TopLevelPath)
	if err != nil {
		return fmt.Errorf("getting stat failed. %w", err)
	}

	if dir.IsDir() {

		args := getArgs(config)

		/* Execute the command */

		cmd = exec.Command("git", args...)
		cmd.Dir = thisRepo.TopLevelPath
		out_bytes, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("getting cmd output failed. %w", err)
		}

		/* Build the output */

		cmdOutput = strings.TrimSpace(string(out_bytes))
		if loggingLevel >= 3 {
			logInfo.Printf("git status command output for %v:\n%s\n", thisRepo, cmdOutput)
		}

		/* Construct ahead / behind symbol */

		regex = regexp.MustCompile(`(?mU)^# branch.ab \+(\d*) \-(\d*)$`) // Capture 2 numbers
		matches = regex.FindStringSubmatch(cmdOutput)
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

		if config.showBranchHead {
			thisRepo.BranchHead = getStringRegex(`(?mU)^# branch.head (.+)$`, cmdOutput)
		}
		if config.showBranchUpstream || config.showUrl {
			thisRepo.BranchUpstream = getStringRegex(`(?mU)^# branch.upstream (.+)$`, cmdOutput)
		}

		return nil
	}
	return errors.New("not a directory")
}

/*
getArgs returns git status command arguments

	'config' rules of the construction

Convenience funtion
*/
func getArgs(config tConfig) []string {
	args := []string{}
	args = append(args,
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
