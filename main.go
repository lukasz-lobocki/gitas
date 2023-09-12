package main

var (
	semVer, commitHash     string
	isGitDirty, isSnapshot string
	goOs, goArch           string
	gitUrl, builtBranch    string
	builtDate, builtBy     string
)

func main() {
	println("build version:", semVer+"+"+goArch+"."+commitHash)
	println()
	println("semVer:", semVer, "commitHash:", commitHash)
	println("isGitDirty:", isGitDirty, "isSnapshot:", isSnapshot)
	println("goOs:", goOs, "goArch:", goArch)
	println("gitUrl:", gitUrl, "builtBranch:", builtBranch)
	println("builtDate:", builtDate, "builtBy:", builtBy)
}
