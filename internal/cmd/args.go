package cmd

var (
	dryRun bool

	checkUpdates bool
	onlyUpdates  bool
	jsonFormat   bool

	listExclude      *[]string
	upgradeExclude   *[]string
	uninstallExclude *[]string
)
