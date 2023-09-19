package cmd

/*
Status' configuration
*/
type tConfig struct {
	sortOrder          *tChoice
	nameShown          *tChoice
	timeFormat         *tChoice
	showUrl            bool
	showCommitTime     bool
	showBranchHead     bool
	showRemoteSyncNeed bool
	showBranchUpstream bool
	showDirty          bool
	showUntracked      bool
	showStash          bool
	lookForSubGits     bool // Not implemented
	emitFormat         *tChoice
}

/*
initChoices sets up Config struct for 'limited choice' flag
*/
func initChoices() {
	config.nameShown = newChoice([]string{"u", "p", "s"}, "u")
	config.sortOrder = newChoice([]string{"t", "n"}, "t")
	config.timeFormat = newChoice([]string{"r", "i"}, "r")
	config.emitFormat = newChoice([]string{"t", "j", "m"}, "t")
}
