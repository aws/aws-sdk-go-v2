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
  tool release <repo path> update-pending`

func usage() {
	fmt.Println(usageText)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error

	switch os.Args[1] {
	case "change":
		err = changeSubcmd(os.Args[2:])
	case "release":
		err = releaseSubcmd(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
	}
}
