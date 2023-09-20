package cmd

// https://mattkrol.me/2020/display-git-status-information-in-your-shell-prompt.html
type tRepo struct {
	TopLevelPath    string `json:"topLevelPath"`    // Full path
	UniqueName      string `json:"uniqueName"`      // Shortest unique path
	TopLevelGroup   string `json:"topLevelGroup"`   // Most significant segment
	ShortName       string `json:"shortName"`       // Least significant segment
	OriginUrl       string `json:"originUrl"`       // github url
	LastCommitTime  string `json:"lastCommitTime"`  // Human-readable
	LastCommitEpoch string `json:"lastCommitEpoch"` // For sorting purposes
	BranchHead      string `json:"branchHead"`
	FetchNeeded     bool   `json:"fetchNeeded"`
	BranchUpstream  string `json:"branchUpstream"`
	Ahead           int    `json:"ahead"`
	Behind          int    `json:"behind"`
	StatusAB        string `json:"statusAB"` // Ahead - behind
	Dirty           bool   `json:"dirty"`
	Untracked       bool   `json:"untracked"`
	Stash           bool   `json:"stash"`
}
