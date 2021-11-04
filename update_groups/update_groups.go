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

const using string = "using System;\nusing UnityEngine;\nusing NaughtyAttributes;\n\n"
const enumHeader string = "public enum UpdateGroup\n{\n"
const classHeader string = "public partial class MonoBehaviourHook : MonoBehaviour\n{\n"
const boxGroup string = "[BoxGroup(\"Groups\")]\n"
const setupGroups string = "void SetupUpdateGroups()\n{\n"
const switchGroup string = "switch (group) {\n"
const updateHeader string = "using UnityEngine;\n\npublic partial class MonoBehaviourHook : MonoBehaviour\n{\nprivate void Update() {\n"
const regUpdate string = "public void RegisterUpdateConsumer(IUpdatable consumer, UpdateGroup group) {\n"
const unregUpdate string = "public void UnregisterUpdateConsumer(IUpdatable consumer, UpdateGroup group) {\n"
const lateUpdateHeader string = "using UnityEngine;\n\npublic partial class MonoBehaviourHook : MonoBehaviour\n{\nprivate void LateUpdate() {\n"
const regLateUpdate string = "public void RegisterLateUpdateConsumer(ILateUpdatable consumer, UpdateGroup group) {\n"
const unregLateUpdate string = "public void UnregisterLateUpdateConsumer(ILateUpdatable consumer, UpdateGroup group) {\n"
const groupsFileCs = "../nothome/Assets/Scripts/MonoBehaviourHookGroups.cs"
const updateFileCs = "../nothome/Assets/Scripts/MonoBehaviourHookUpdate.cs"
const lateUpdateFileCs = "../nothome/Assets/Scripts/MonoBehaviourHooklateUpdate.cs"

func UpdateGroups() {
	groupsFile, err := ioutil.ReadFile("./data/updategroups.yaml")
	if err != nil {
		fmt.Println("Failed opening updategroups.yaml")
	}

	groupList := &UpdateGroupList{}
	err = yaml.Unmarshal(groupsFile, groupList)
	if err != nil {
		fmt.Println("Failed unmarshalling updategroups.yaml")
	}

	err = CreateGroupsFile(*groupList)
	if err != nil {
		fmt.Println("Failed Updating MonoBehaviourHookGroups.cs")
	}

	err = CreateUpdateFile(*groupList)
	if err != nil {
		fmt.Println("Failed Updating MonoBehaviourHookUpdate.cs")
	}

	err = CreateLateUpdateFile(*groupList)
	if err != nil {
		fmt.Println("Failed Updating MonoBehaviourHookLateUpdate.cs")
	}
}

func CreateGroupsFile(groups UpdateGroupList) error {
	os.Truncate(groupsFileCs, 0)
	file, err := os.OpenFile(groupsFileCs, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
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
	for _, g := range groups.Groups {
		result += boxGroup
		result += fmt.Sprintf("public bool %sActive;\n", g)
		result += fmt.Sprintf("public event Action On%sUpdate;\n", g)
		result += fmt.Sprintf("public event Action On%sLateUpdate;\n\n", g)
	}

	result += setupGroups
	for _, g := range groups.Groups {
		result += fmt.Sprintf("UpdateGroupsActive.Add(UpdateGroup.%s, %sActive);\n", g, g)
	}
	result += "}\n}"

	_, _ = datawriter.WriteString(result)
	datawriter.Flush()
	file.Close()

	return nil
}

func CreateUpdateFile(groups UpdateGroupList) error {
	os.Truncate(updateFileCs, 0)
	file, err := os.OpenFile(updateFileCs, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	datawriter := bufio.NewWriter(file)

	var result string

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
	result += "default:\nbreak;\n}\n}\n}"

	_, _ = datawriter.WriteString(result)
	datawriter.Flush()
	file.Close()

	return nil
}

func CreateLateUpdateFile(groups UpdateGroupList) error {
	os.Truncate(lateUpdateFileCs, 0)
	file, err := os.OpenFile(lateUpdateFileCs, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	datawriter := bufio.NewWriter(file)

	var result string

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

	return nil
}
