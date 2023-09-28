package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

/*
getReposDictionary returns a slice of Repos found under the 'dirName'

	'dirName' path to be searched
	'config' rules of the search
*/
func getReposDictionary(dirName string, config tConfig) ([]tRepo, error) {

	var (
		allRepos    []tRepo
		thisSpinner *spinner.Spinner = nil
	)

	/* Get repos slice */

	gitsSlice, err := findRepos(dirName, config.lookForSubGits)
	if err != nil {
		return nil, fmt.Errorf("finding repos failed. %w", err)
	}
	if loggingLevel >= 2 {
		logInfo.Printf("%d repos found.", len(gitsSlice))
	}

	/* Get common path prefix for all the repos */

	commonPrefix := commonPrefix(os.PathSeparator, gitsSlice) + string(os.PathSeparator)
	if loggingLevel >= 2 {
		logInfo.Printf("common prefix: %s", commonPrefix)
	}

	/* Main loop */

	thisSpinner = spinner.New(spinner.CharSets[14], time.Duration(SPINNER_MS)*time.Millisecond, spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Retrieving status of repositories\n"))
	thisSpinner.Start() // Starting spinner to show visual work

	for _, thisGit := range gitsSlice {

		var thisRepo tRepo

		/* Putting full path */

		thisRepo.TopLevelPath = thisGit

		/* Setting names */

		// Least significant segment
		thisRepo.ShortName = thisRepo.TopLevelPath[strings.LastIndex(thisRepo.TopLevelPath, string(os.PathSeparator))+1:]
		// Removing commonPrefix
		thisRepo.UniqueName = strings.ReplaceAll(thisRepo.TopLevelPath, commonPrefix, "")
		// Most significant segment from UniqueName
		thisRepo.TopLevelGroup = strings.Split(thisRepo.UniqueName, string(os.PathSeparator))[0]

		/* Edge case, when there is only one repo found */

		if len(gitsSlice) == 1 {
			thisRepo.UniqueName = thisRepo.ShortName
		}

		thisSpinner.Suffix = " " + thisRepo.ShortName

		/* Get most part of repo's status */

		if err := getRepoStatus(&thisRepo, config); err != nil {
			return nil, fmt.Errorf("getting repos status failed. %w", err)
		}

		/* Get fetch needed */

		if config.showFetchNeeded && len(thisRepo.BranchUpstream) > 0 {
			if err = getFetchNeeded(&thisRepo); err != nil {
				return nil, fmt.Errorf("getting remote sync need failed. %w", err)
			}
		}

		/* Get repo's url */

		if config.showUrl && len(thisRepo.BranchUpstream) > 0 {
			if err = getOriginUrl(&thisRepo); err != nil {
				return nil, fmt.Errorf("getting origin url failed. %w", err)
			}
		}

		/* Get repo's time */

		if config.showCommitTime {
			if err = getLastCommitTime(&thisRepo, false, config); err != nil {
				return nil, fmt.Errorf("getting last commit time failed. %w", err)
			}
		}

		/* Get repo's epoch */

		if config.sortOrder.Value == "t" {
			if err = getLastCommitTime(&thisRepo, true, config); err != nil {
				return nil, fmt.Errorf("getting last commit epoch failed. %w", err)
			}
		}

		/* Append thisRepo to allRepos */

		allRepos = append(allRepos, thisRepo)

	}

	thisSpinner.Stop()

	return allRepos, nil
}

/*
findRepos returns the slice of top-level gir repo paths

	'dirName' path to search and its children
	'lookForSubGits' not implemented
*/
func findRepos(dirName string, lookForSubGits bool) ([]string, error) {

	var (
		thisResult  []string
		thisSpinner *spinner.Spinner = nil
	)

	thisSpinner = spinner.New(spinner.CharSets[14], time.Duration(SPINNER_MS)*time.Millisecond, spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Finding repositories\n"))
	thisSpinner.Start() // Starting spinner to show visual work

	fsys := os.DirFS(dirName)

	err := fs.WalkDir(
		fsys,
		".",
		func(thisPath string, thisDir fs.DirEntry, err error) error {

			var thisFullPath = filepath.Join(dirName, thisPath)

			if err != nil {
				return fmt.Errorf("findRepos: WalkDirFunc: error accessing path %s", thisFullPath)
			}

			/* Hop over the files */

			if !thisDir.IsDir() {
				// Skipping
				return nil
			}

			/* Hop over non-git dirs */

			var isInGit bool

			if isInGit, err = isInGitWorkTree(thisFullPath); err != nil {
				return fmt.Errorf("findRepos: isInGitWorkTree failed. %w", err)
			}

			if !isInGit {
				// Skipping
				return nil
			}

			/* Get top level git path */

			var gitTopLevel string

			if gitTopLevel, err = getGitTopLevel(thisFullPath); err != nil {
				return fmt.Errorf("findRepos: gitTopLevel failed. %w", err)
			}

			thisResult = append(thisResult, gitTopLevel)
			if loggingLevel >= 3 {
				logInfo.Printf("findRepos: appending: %s", gitTopLevel)
			}

			/* Hop over every dir beneath the one */

			if !lookForSubGits {
				if loggingLevel >= 3 {
					logInfo.Println("Skipping dir")
				}
				// Skipping
				return filepath.SkipDir
			}

			return nil
		},
	)
	if err != nil {
		return []string{}, fmt.Errorf("findRepos: WalkDir failed. %w", err)
	}

	if loggingLevel >= 2 {
		logInfo.Printf("findRepos: thisResult: %+v\n", thisResult)
	}

	thisSpinner.Stop()

	return thisResult, nil
}

/*
isInGitWorkTree returns if path is within git work tree

	'dirName' path to be checked
*/
func isInGitWorkTree(dirName string) (bool, error) {

	/* Get os object */
	if loggingLevel >= 3 {
		logInfo.Printf("isInGitWorkTree: dirName: %s", dirName)
	}

	dir, err := os.Stat(dirName)
	if err != nil {
		return false, fmt.Errorf("isInGitWorkTree: getting stat failed. %w", err)
	}

	if !dir.IsDir() {
		if loggingLevel >= 3 {
			logInfo.Printf("isInGitWorkTree: not a dir: %s", dirName)
		}
		return false, nil
	}

	/* Execute command */

	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dirName
	out_bytes, err := cmd.Output()
	if err != nil {
		return false, nil // Not in git work tree anyway
	}

	return strings.TrimSpace(string(out_bytes)) == "true", nil

}

/*
getGitTopLevel returns top-level path of git repo

	'dirName' path to repo or any of its child directories
*/
func getGitTopLevel(dirName string) (string, error) {

	dir, err := os.Stat(dirName)
	if err != nil {
		return "", fmt.Errorf("getting stat failed. %w", err)
	}

	if dir.IsDir() {

		/* Execute command */

		cmd := exec.Command("git", "rev-parse", "--show-toplevel")
		cmd.Dir = dirName
		out_bytes, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("getting cmd output failed. %w", err)
		}

		return strings.TrimSpace(string(out_bytes)), nil
	}

	return "", errors.New("not a directory")
}
