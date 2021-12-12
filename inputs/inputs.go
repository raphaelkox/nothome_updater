package inputs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type InputGroupList struct {
	Groups []InputGroup `yaml:"groups"`
}

type InputGroup struct {
	Name   string
	Inputs []InputData `yaml:"inputs"`
}

type InputData struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

const using string = "using System;\nusing UnityEngine;\n\n"
const classHeader string = "public static partial class InputManager {\n"

func UpdateInputGroups() {
	inputsFile, err := ioutil.ReadFile("./data/inputs.yaml")
	if err != nil {
		fmt.Println("Failed opening inputs.yaml")
	}

	inputGroupList := &InputGroupList{}
	err = yaml.Unmarshal(inputsFile, inputGroupList)
	if err != nil {
		fmt.Println("Failed unmarshalling inputs.yaml")
	}

	filePath := "../nothome/Assets/Scripts/InputSystem/InputManager.cs"
	os.Truncate(filePath, 0)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed opening output file: InputManager.cs")
		return
	}
	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString("namespace NH.Input\n{\n")
	_, _ = datawriter.WriteString("using Core;\n\n")
	_, _ = datawriter.WriteString("public static partial class InputManager\n{\n")
	_, _ = datawriter.WriteString("static bool Registered = false;\n\n")
	_, _ = datawriter.WriteString("public static void Setup() {\nif (!Registered) {\n")

	for _, group := range inputGroupList.Groups {
		_, _ = datawriter.WriteString(fmt.Sprintf("UpdateGroupsManager.On%sInputUpdate += %sInputGroup.OnUpdate;\n", group.Name, group.Name))
		updateInputs(&group)
	}

	_, _ = datawriter.WriteString("Registered = true;\n}\n}\n}\n}")

	datawriter.Flush()
	file.Close()
}

func updateInputs(inputGroup *InputGroup) {
	updateInputGroup(inputGroup)
	updateInputConsumer(inputGroup)
}

func updateInputGroup(inputGroup *InputGroup) {
	name := inputGroup.Name
	filePath := fmt.Sprintf("../nothome/Assets/Scripts/InputSystem/%sInputGroup.cs", name)
	os.Truncate(filePath, 0)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed opening output file: %sInputGroup.cs\n", name)
		return
	}
	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString(using)
	_, _ = datawriter.WriteString("namespace NH.Input\n{")
	_, _ = datawriter.WriteString("using Interfaces;\n")
	_, _ = datawriter.WriteString("using Input = UnityEngine.Input;\n\n")
	_, _ = datawriter.WriteString(fmt.Sprintf("public static class %sInputGroup \n{\n", name))

	_, _ = datawriter.WriteString(createStateStruct(inputGroup))
	_, _ = datawriter.WriteString(fmt.Sprintf("static %sInputState currentState;\n", inputGroup.Name))
	_, _ = datawriter.WriteString(fmt.Sprintf("static %sInputState prevState;\n\n", inputGroup.Name))

	_, _ = datawriter.WriteString(fmt.Sprintf("public static %sInputState CurrentState {\nget { return currentState; }\n}\n", inputGroup.Name))
	_, _ = datawriter.WriteString(fmt.Sprintf("public static %sInputState PrevState {\nget { return prevState; }\n}\n", inputGroup.Name))

	_, _ = datawriter.WriteString(createUpdateMethod(inputGroup.Inputs))

	for _, input := range inputGroup.Inputs {
		_, _ = datawriter.WriteString(createInput(input.Name, input.Key))
	}

	_, _ = datawriter.WriteString("}\n}")

	datawriter.Flush()
	file.Close()
}

func updateInputConsumer(inputGroup *InputGroup) {
	name := inputGroup.Name
	filePath := fmt.Sprintf("../nothome/Assets/Scripts/InputSystem/%sInputConsumer.cs", name)
	os.Truncate(filePath, 0)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed opening output file: Input%sConsumer.cs\n", name)
		return
	}
	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString("using UnityEngine;\n")
	_, _ = datawriter.WriteString(fmt.Sprintf("public class %sInputConsumer : MonoBehaviour\n{\n", name))
	_, _ = datawriter.WriteString("[SerializeField] CustomBehaviour consumer;\n\n")

	for _, input := range inputGroup.Inputs {
		_, _ = datawriter.WriteString(fmt.Sprintf("[SerializeField] bool %s;\n", input.Name))
	}

	_, _ = datawriter.WriteString("\nprivate void OnEnable() {\n")

	for _, input := range inputGroup.Inputs {
		_, _ = datawriter.WriteString(fmt.Sprintf("if (%s) {\n%sInputGroup.Register%sConsumer(consumer);\n}\n", input.Name, inputGroup.Name, input.Name))
	}

	_, _ = datawriter.WriteString("}\n}")

	datawriter.Flush()
	file.Close()
}

func createStateStruct(inputGroup *InputGroup) string {
	var result string

	result += fmt.Sprintf("public struct %sInputState {\n", inputGroup.Name)
	for _, input := range inputGroup.Inputs {
		result += fmt.Sprintf("public bool %s;\n", strings.ToLower(input.Name))
	}

	result += "}\n\n"

	return result
}

func createUpdateMethod(inputGroup []InputData) string {
	var result string

	result += fmt.Sprintf("public static void OnUpdate() {\n")
	result += fmt.Sprintf("prevState = currentState;\n\n")

	for _, input := range inputGroup {
		result += fmt.Sprintf("Update%sInput();\n", input.Name)
	}

	result += "}\n\n"

	return result
}

func createInput(name string, key string) string {
	var result string

	result += fmt.Sprintf("\n#region %s\n", name)
	result += fmt.Sprintf("public static event Action On%sInputDown;\n", name)
	result += fmt.Sprintf("public static event Action On%sInputUp;\n", name)
	result += fmt.Sprintf("static void Update%sInput() {\n", name)
	result += fmt.Sprintf("currentState.%s = Input.GetKey(KeyCode.%s);\n", strings.ToLower(name), key)
	result += fmt.Sprintf("if (Input.GetKeyDown(KeyCode.%s)) {\n", key)
	result += fmt.Sprintf("On%sInputDown?.Invoke();\n}\n", name)
	result += fmt.Sprintf("if (Input.GetKeyUp(KeyCode.%s)) {\n", key)
	result += fmt.Sprintf("On%sInputUp?.Invoke();\n}\n}\n", name)
	result += fmt.Sprintf("public static void Register%sConsumer(IRegistrable consumer) {\n", name)
	result += fmt.Sprintf("consumer.OnDisableEvent += Unregister%sConsumer;\n", name)
	result += fmt.Sprintf("consumer.OnDestroyEvent += Unregister%sConsumer;\n\n", name)
	result += fmt.Sprintf("I%sInputConsumer %sConsumer = consumer as I%sInputConsumer;\n", name, name, name)
	result += fmt.Sprintf("On%sInputDown += %sConsumer.Handle%sInputDown;\n", name, name, name)
	result += fmt.Sprintf("On%sInputUp += %sConsumer.Handle%sInputUp;\n}\n", name, name, name)
	result += fmt.Sprintf("public static void Unregister%sConsumer(IRegistrable consumer) {\n", name)
	result += fmt.Sprintf("consumer.OnDisableEvent -= Unregister%sConsumer;\n", name)
	result += fmt.Sprintf("consumer.OnDestroyEvent -= Unregister%sConsumer;\n\n", name)
	result += fmt.Sprintf("I%sInputConsumer %sConsumer = consumer as I%sInputConsumer;\n", name, name, name)
	result += fmt.Sprintf("On%sInputDown -= %sConsumer.Handle%sInputDown;\n", name, name, name)
	result += fmt.Sprintf("On%sInputUp -= %sConsumer.Handle%sInputUp;\n}\n", name, name, name)
	result += fmt.Sprintf("public interface I%sInputConsumer {\n", name)
	result += fmt.Sprintf("public void Handle%sInputDown();\n", name)
	result += fmt.Sprintf("public void Handle%sInputUp();\n}\n", name)
	result += fmt.Sprintf("#endregion\n")

	return result
}
