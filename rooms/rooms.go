package rooms

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type RoomList struct {
	Rooms []RoomData `yaml:"rooms"`
}

type RoomData struct {
	Name        string   `yaml:"name"`
	SpawnPoints []string `yaml:"points"`
}

const using string = "using System;\nusing System.Collections.Generic;\nusing UnityEngine;\n\n"
const classHeader string = "public static class RoomSystem\n{\n"

func UpdateRooms() {
	roomsFile, err := ioutil.ReadFile("./data/rooms.yaml")
	if err != nil {
		fmt.Println("Failed opening rooms.yaml")
	}

	roomList := &RoomList{}
	err = yaml.Unmarshal(roomsFile, roomList)
	if err != nil {
		fmt.Println("Failed unmarshalling rooms.yaml")
	}

	filePath := "../nothome/Assets/Scripts/RoomSystem.cs"
	os.Truncate(filePath, 0)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString(using)
	_, _ = datawriter.WriteString(classHeader)
	_, _ = datawriter.WriteString("public static string NextRoomSpawn = \"RIGHT\";\n\n")
	_, _ = datawriter.WriteString("public static Dictionary<string, System.Type> SpawnPointLists = new Dictionary<string, System.Type>() {\n")

	for _, r := range roomList.Rooms {
		_, _ = datawriter.WriteString(fmt.Sprintf("{\"%s\", typeof(%s_spawnpoints) },\n", r.Name, r.Name))
	}
	_, _ = datawriter.WriteString("};\n")

	for _, r := range roomList.Rooms {
		_, _ = datawriter.WriteString(fmt.Sprintf("\npublic enum %s_spawnpoints\n{\n", r.Name))
		for _, sp := range r.SpawnPoints {
			_, _ = datawriter.WriteString(fmt.Sprintf("%s,\n", sp))
		}
		_, _ = datawriter.WriteString("}\n")
	}

	_, _ = datawriter.WriteString("}")

	datawriter.Flush()
	file.Close()
}
