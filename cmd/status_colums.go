package cmd

import (
	"strings"

	"github.com/fatih/color"
)

type tColumn struct {
	isShown         func(tConfig) bool
	title           func(tConfig) string
	titleColor      color.Attribute
	contentSource   func(tConfig, tRepo) string
	contentColor    func(tRepo) color.Attribute
	contentAlignMD  int
	contentEscapeMD bool
}

/*
getClickable returns linkText string that is clickable in GNOME and spawns url
*/
func getClickable(linkText string, url string) string {
	return "\033]8;;" + url + "\a" + linkText + "\033]8;;\a"
}

/*
getColumns defines look and content of table's emitted columns
*/
func getColumns() []tColumn {

	var thisColumns []tColumn

	thisColumns = append(thisColumns,

		tColumn{ // Name
			isShown: func(_ tConfig) bool { return true }, // Always shown
			title: func(tc tConfig) string {
				switch tc.nameShown.Value { // Title differs by config
				case "p":
					return "Top-level path"
				case "s":
					return "Short"
				case "u":
					return "Unique name"
				}
				return ""
			},
			titleColor: color.Bold,

			contentSource: func(tc tConfig, tr tRepo) string {
				switch tc.nameShown.Value { // Content differs by config
				case "p":
					return getClickable(tr.TopLevelPath, "file:///"+tr.TopLevelPath)
				case "s":
					return getClickable(tr.ShortName, "file:///"+tr.TopLevelPath)
				case "u":
					return getClickable(tr.UniqueName, "file:///"+tr.TopLevelPath)
				}
				return ""
			},
			contentColor:    func(_ tRepo) color.Attribute { return color.FgHiYellow }, // Static color
			contentAlignMD:  ALIGN_LEFT,
			contentEscapeMD: true,
		},

		tColumn{ // showCommitTime
			isShown:    func(tc tConfig) bool { return tc.showCommitTime },
			title:      func(_ tConfig) string { return "Last commit" }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return tr.LastCommitTime },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgHiBlack }, // Static color
			contentAlignMD:  ALIGN_LEFT,
			contentEscapeMD: true,
		},

		tColumn{ // showBranchHead
			isShown:    func(tc tConfig) bool { return tc.showBranchHead },
			title:      func(_ tConfig) string { return "Branch head" }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return tr.BranchHead },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgHiBlue }, // Static color
			contentAlignMD:  ALIGN_LEFT,
			contentEscapeMD: true,
		},
		tColumn{ // showFetchNeeded
			isShown:    func(tc tConfig) bool { return tc.showFetchNeeded },
			title:      func(_ tConfig) string { return "Q" }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return parseBool(tr.FetchNeeded, FETCH_NEEDED_SYMBOL) },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgHiCyan }, // Static color
			contentAlignMD:  ALIGN_CENTER,
			contentEscapeMD: false,
		},
		tColumn{ // showBranchUpstream
			isShown:    func(tc tConfig) bool { return tc.showBranchUpstream },
			title:      func(_ tConfig) string { return "Branch remote" }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return tr.BranchUpstream },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgHiBlue }, // Static color
			contentAlignMD:  ALIGN_LEFT,
			contentEscapeMD: true,
		},

		tColumn{ // showUrl
			isShown:    func(tc tConfig) bool { return tc.showUrl },
			title:      func(_ tConfig) string { return "Url" }, // Static title
			titleColor: color.Bold,

			contentSource: func(_ tConfig, tr tRepo) string {
				return getClickable(tr.OriginUrl, strings.ReplaceAll(tr.OriginUrl, "ssh://git@", "https://"))
				/* return strings.ReplaceAll(
					tr.OriginUrl, "git@github.com:", "ssh@https://github.com/", // To provide clickable text in the output
				) */
			},
			contentColor:    func(_ tRepo) color.Attribute { return color.FgWhite }, // Static color
			contentAlignMD:  ALIGN_LEFT,
			contentEscapeMD: false,
		},

		tColumn{ // Ahead / behind
			isShown:    func(tc tConfig) bool { return true },               // Always shown
			title:      func(_ tConfig) string { return PUSH_FETCH_SYMBOL }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return getThisABSymbol()[tr.StatusAB] },
			contentColor:    func(tr tRepo) color.Attribute { return getThisABColor()[tr.StatusAB] }, // Dynamic color
			contentAlignMD:  ALIGN_CENTER,
			contentEscapeMD: false,
		},

		tColumn{ // showDirty
			isShown:    func(tc tConfig) bool { return tc.showDirty },
			title:      func(_ tConfig) string { return "D" }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return parseBool(tr.Dirty, DIRTY_SYMBOL) },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgCyan }, // Static color
			contentAlignMD:  ALIGN_CENTER,
			contentEscapeMD: false,
		},

		tColumn{ // showUntracked
			isShown:    func(tc tConfig) bool { return tc.showUntracked },
			title:      func(_ tConfig) string { return "U" }, // Static title
			titleColor: color.Bold,

			contentSource:   func(_ tConfig, tr tRepo) string { return parseBool(tr.Untracked, UNTRACKED_SYMBOL) },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgRed }, // Static color
			contentAlignMD:  ALIGN_CENTER,
			contentEscapeMD: false,
		},

		tColumn{ // showStash
			isShown: func(tc tConfig) bool { return tc.showStash },
			title:   func(_ tConfig) string { return "S" }, // Static title

			titleColor:      color.Bold,
			contentSource:   func(_ tConfig, tr tRepo) string { return parseBool(tr.Stash, STASH_SYMBOL) },
			contentColor:    func(_ tRepo) color.Attribute { return color.FgYellow }, // Static color
			contentAlignMD:  ALIGN_CENTER,
			contentEscapeMD: false,
		},
	)

	return thisColumns
}

/*
Ahead / behind symbols

⇅ ↑ ↓ ⬍ ⬆ ⬇ ⭡ ⭣ ⮁ ⁕ ⁂ ⨹ ⊡ ⧆ ⊙ ⊛ ⨻ ⊠ ⊞ ⎗ ⮧ ⮯ ⮧ ⊝ ⊗ ⊜ ↥ ↧ ↭ ↹ ↗ ↘ ↯ ◥ ◢ ◹ ◿ ⤓
*/
const (
	SYNCED_SYMBOL       string = "✓"
	REMOTE_AHEAD_SYMBOL string = "↘"
	LOCAL_AHEAD_SYMBOL  string = "↗"
	DIVERGED_SYMBOL     string = "↹"
	PUSH_FETCH_SYMBOL   string = "⇅"
)

const (
	SYNCED_CHAR       string = "synced"
	REMOTE_AHEAD_CHAR string = "ready for merge"
	LOCAL_AHEAD_CHAR  string = "ready for push"
	DIVERGED_CHAR     string = "diverged"
)

/*
getThisABSymbol maps given ahead / behind status string to appropriate symbol
*/
func getThisABSymbol() map[string]string {
	return map[string]string{
		SYNCED_CHAR:       SYNCED_SYMBOL,
		REMOTE_AHEAD_CHAR: REMOTE_AHEAD_SYMBOL,
		LOCAL_AHEAD_CHAR:  LOCAL_AHEAD_SYMBOL,
		DIVERGED_CHAR:     DIVERGED_SYMBOL,
	}
}

/*
getThisABColor maps given ahead / behind status string to appropriate color
*/
func getThisABColor() map[string]color.Attribute {
	return map[string]color.Attribute{
		SYNCED_CHAR:       color.FgHiGreen,
		REMOTE_AHEAD_CHAR: color.FgHiCyan,
		LOCAL_AHEAD_CHAR:  color.FgHiMagenta,
		DIVERGED_CHAR:     color.FgHiRed,
	}
}

const (
	ALIGN_LEFT = iota
	ALIGN_CENTER
	ALING_RIGHT
)

/*
getThisAlignChar amps given alignment to appropriate markdown string to be used in header separator
*/
func getThisAlignChar() map[int]string {
	return map[int]string{
		ALIGN_LEFT:   `:-`,
		ALIGN_CENTER: `:-:`,
		ALING_RIGHT:  `-:`,
	}
}

/*
Other symbols
*/
const (
	DIRTY_SYMBOL        string = "⊛"
	UNTRACKED_SYMBOL    string = "⊗"
	STASH_SYMBOL        string = "⊜"
	FETCH_NEEDED_SYMBOL string = "↯" // Ready for fetch
)
