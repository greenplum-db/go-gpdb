package argParser

import (
	"fmt"
	"os"
)

// Program Usage.
func ShowHelp() {
	fmt.Print(`
Usage: gpdb [COMMAND]
COMMAND:
	download        Download software from PivNet
	install         Install GPDB software on the host
	remove          Remove a perticular Installation
	env             Show all the environment of installed version
	version         Show version of the script
	help            Show help
`)
	os.Exit(0)
}
