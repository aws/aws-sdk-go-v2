package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("USAGE:\n  tool change add [-module=<module>]\n  tool change ls [-module=<module>]\n  " +
		"tool change modify <change_id>\n  tool change rm <change_id>\n  tool release update-pending\n")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "change":
		changeSubcmd(os.Args[2:])
	case "release":
		releaseSubcmd(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}
