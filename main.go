package main

import (
	"fmt"
	"os"
	"updater/inputs"
	updategroups "updater/update_groups"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}

	switch os.Args[1] {
	case "inputs":
		inputs.UpdateInputGroups()
	case "update_groups":
		updategroups.UpdateGroups()
	default:
		fmt.Println("Invalid argument")
	}
}
