package main

import (
	"fmt"
	"os"
	"updater/input"
	"updater/rooms"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}

	switch os.Args[1] {
	case "input":
		input.Updateinput()
	case "rooms":
		rooms.UpdateRooms()
	default:
		fmt.Println("Invalid argument")
	}
}
