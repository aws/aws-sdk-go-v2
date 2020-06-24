package main

import (
	"fmt"
	"os"
)

const usageText = `USAGE:
  tool change add [-module=<module>]
  tool change ls [-module=<module>]
  tool change modify <change_id>
  tool change rm <change_id>
  tool release update-pending`

func usage() {
	fmt.Println(usageText)
	os.Exit(1)
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
