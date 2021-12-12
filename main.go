package main

import (
	"fmt"
	"os"
	"updater/inputs"
	"updater/rooms"
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
	case "rooms":
		rooms.UpdateRooms()
	case "update_groups":
		updategroups.UpdateGroups()
	default:
		fmt.Println("Invalid argument")
	}
}
