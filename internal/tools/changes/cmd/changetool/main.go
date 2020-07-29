package main

import (
	"fmt"
	"os"
)

const usageText = `USAGE:
  tool change add [-module=<module>] [-type=<type>] [-description=<description>]
  tool change ls [-module=<module>]
  tool change modify <change_id>
  tool change rm <change_id>
  tool release <repo path> update-pending`

func usage() {
	changeUsage()
	fmt.Println()
	releaseUsage()
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	var err error

	switch os.Args[1] {
	case "change":
		err = changeSubcmd(os.Args[2:])
	case "release":
		err = releaseSubcmd(os.Args[2:])
	default:
		usage()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
