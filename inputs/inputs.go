package inputs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type InputList struct {
	Inputs []InputData `yaml:"inputs"`
}

type InputData struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

const using string = "using System;\nusing UnityEngine;\n\n"
const classHeader string = "public static partial class InputManager {\n"
const structHeader string = "public struct InputState {\n"

func Updateinput() {
	inputsFile, err := ioutil.ReadFile("./data/inputs.yaml")
	if err != nil {
		fmt.Println("Failed opening inputs.yaml")
	}

	inputList := &InputList{}
	err = yaml.Unmarshal(inputsFile, inputList)
	if err != nil {
		fmt.Println("Failed unmarshalling inputs.yaml")
	}

	filePath := "../nothome/Assets/Scripts/InputManagerMethods.cs"
	os.Truncate(filePath, 0)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	datawriter := bufio.NewWriter(file)

	_, _ = datawriter.WriteString(using)
	_, _ = datawriter.WriteString(classHeader)
	_, _ = datawriter.WriteString(CreateStateStruct(inputList.Inputs))
	_, _ = datawriter.WriteString(CreateUpdateMethod(inputList.Inputs))
	_, _ = datawriter.WriteString(CreateDirectionInput())

	for _, input := range inputList.Inputs {
		_, _ = datawriter.WriteString(CreateButtonInput(input.Name, input.Key))
	}

	_, _ = datawriter.WriteString("}")

	datawriter.Flush()
	file.Close()
}

func CreateStateStruct(inputList []InputData) string {
	var result string

	result += structHeader
	result += "public Vector2 direction;\n"

	for _, input := range inputList {
		result += fmt.Sprintf("public bool %s;\n", strings.ToLower(input.Name))
	}

	result += "}\n\n"

	return result
}

func CreateUpdateMethod(inputList []InputData) string {
	var result string

	result += fmt.Sprintf("public static void OnUpdate() {\n")
	result += fmt.Sprintf("prevState = currentState;\n\n")
	result += fmt.Sprintf("UpdateDirectionInput();\n")

	for _, input := range inputList {
		result += fmt.Sprintf("Update%sInput();\n", input.Name)
	}

	result += "}\n\n"

	return result
}

func CreateDirectionInput() string {
	var result string

	result += `#region Direction
    public static event Action<Vector2> OnDirInput;
    static void UpdateDirectionInput() {
        currentState.direction = new Vector2(Input.GetAxisRaw("Horizontal"), Input.GetAxisRaw("Vertical"));
        if (currentState.direction != prevState.direction) {
            OnDirInput?.Invoke(currentState.direction);
        }
    }
    public static void RegisterDirectionConsumer(IDirectionInputConsumer consumer) {
        OnDirInput += consumer.HandleDirectionInput;
    }
    public static void UnregisterDirectionConsumer(IDirectionInputConsumer consumer) {
        OnDirInput -= consumer.HandleDirectionInput;
    }
    public interface IDirectionInputConsumer
    {
        public void HandleDirectionInput(Vector2 direction);
    }
    #endregion`

	result += "\n"

	return result
}

func CreateButtonInput(name string, key string) string {
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
	result += fmt.Sprintf("public static void Register%sConsumer(I%sInputConsumer consumer) {\n", name, name)
	result += fmt.Sprintf("On%sInputDown += consumer.Handle%sInputDown;\n", name, name)
	result += fmt.Sprintf("On%sInputUp += consumer.Handle%sInputUp;\n}\n", name, name)
	result += fmt.Sprintf("public static void Unregister%sConsumer(I%sInputConsumer consumer) {\n", name, name)
	result += fmt.Sprintf("On%sInputDown -= consumer.Handle%sInputDown;\n", name, name)
	result += fmt.Sprintf("On%sInputUp -= consumer.Handle%sInputUp;\n}\n", name, name)
	result += fmt.Sprintf("public interface I%sInputConsumer {\n", name)
	result += fmt.Sprintf("public void Handle%sInputDown();\n", name)
	result += fmt.Sprintf("public void Handle%sInputUp();\n}\n", name)
	result += fmt.Sprintf("#endregion\n")

	return result
}
