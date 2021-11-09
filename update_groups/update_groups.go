package updategroups

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type UpdateGroupList struct {
	Groups []string `yaml:"groups"`
}

const using string = "using System;\nusing System.Collections.Generic;\n\n"
const enumHeader string = "public enum UpdateGroup\n{\n"
const classHeader string = "public static class UpdateGroupsManager\n{\n"
const dictionaryHeader = "private static Dictionary<UpdateGroup, bool> UpdateGroupsActive = new Dictionary<UpdateGroup, bool>() {\n"
const groupActiveFunc = "public static bool IsGroupActive(UpdateGroup group) {\nreturn UpdateGroupsActive[group];\n}\n"
const setGroupsStateFunc = "public static void SetUpdateGroupState(UpdateGroup group, bool state) {\nUpdateGroupsActive[group] = state;\n}\n\n"
const switchGroup string = "switch (group) {\n"
const updateHeader string = "\npublic static void OnUpdate() {\n"
const regUpdate string = "public static void RegisterUpdateConsumer(IUpdatable consumer, UpdateGroup group) {\n"
const unregUpdate string = "public static void UnregisterUpdateConsumer(IUpdatable consumer, UpdateGroup group) {\n"
const lateUpdateHeader string = "\npublic static void OnLateUpdate() {\n"
const regLateUpdate string = "public static void RegisterLateUpdateConsumer(ILateUpdatable consumer, UpdateGroup group) {\n"
const unregLateUpdate string = "public static void UnregisterLateUpdateConsumer(ILateUpdatable consumer, UpdateGroup group) {\n"

const groupsFileCs = "../nothome/Assets/Scripts/Core/UpdateGroupsManager.cs"

func UpdateGroups() {
	groupsFile, err := ioutil.ReadFile("./data/updategroups.yaml")
	if err != nil {
		fmt.Println("Failed opening updategroups.yaml")
	}

	groups := &UpdateGroupList{}
	err = yaml.Unmarshal(groupsFile, groups)
	if err != nil {
		fmt.Println("Failed unmarshalling updategroups.yaml")
	}

	os.Truncate(groupsFileCs, 0)
	file, err := os.OpenFile(groupsFileCs, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed opening updategroups.yaml")
	}
	datawriter := bufio.NewWriter(file)

	var result string

	result += using
	result += enumHeader

	for _, g := range groups.Groups {
		result += fmt.Sprintf("%s,\n", g)
	}
	result += "}\n\n"

	result += classHeader
	result += dictionaryHeader

	for _, g := range groups.Groups {
		result += fmt.Sprintf("{UpdateGroup.%s, true},\n", g)
	}

	result += "};\n\n"

	result += groupActiveFunc
	result += setGroupsStateFunc

	for _, g := range groups.Groups {
		result += fmt.Sprintf("public static event Action On%sUpdate;\n", g)
	}

	result += updateHeader

	for _, g := range groups.Groups {
		result += fmt.Sprintf("if (UpdateGroupsActive[UpdateGroup.%s]) {\nOn%sUpdate?.Invoke();\n}\n", g, g)
	}
	result += "}\n\n"

	result += regUpdate
	result += switchGroup
	for _, g := range groups.Groups {
		result += fmt.Sprintf("case UpdateGroup.%s:\nOn%sUpdate += consumer.OnUpdate;\nbreak;\n", g, g)
	}
	result += "default:\nbreak;\n}\n}\n\n"

	result += unregUpdate
	result += switchGroup
	for _, g := range groups.Groups {
		result += fmt.Sprintf("case UpdateGroup.%s:\nOn%sUpdate -= consumer.OnUpdate;\nbreak;\n", g, g)
	}
	result += "default:\nbreak;\n}\n}\n\n"

	for _, g := range groups.Groups {
		result += fmt.Sprintf("public static event Action On%sLateUpdate;\n", g)
	}

	result += lateUpdateHeader

	for _, g := range groups.Groups {
		result += fmt.Sprintf("if (UpdateGroupsActive[UpdateGroup.%s]) {\nOn%sLateUpdate?.Invoke();\n}\n", g, g)
	}
	result += "}\n\n"

	result += regLateUpdate
	result += switchGroup
	for _, g := range groups.Groups {
		result += fmt.Sprintf("case UpdateGroup.%s:\nOn%sLateUpdate += consumer.OnLateUpdate;\nbreak;\n", g, g)
	}
	result += "default:\nbreak;\n}\n}\n\n"

	result += unregLateUpdate
	result += switchGroup
	for _, g := range groups.Groups {
		result += fmt.Sprintf("case UpdateGroup.%s:\nOn%sLateUpdate -= consumer.OnLateUpdate;\nbreak;\n", g, g)
	}
	result += "default:\nbreak;\n}\n}\n}"

	_, _ = datawriter.WriteString(result)
	datawriter.Flush()
	file.Close()
}
