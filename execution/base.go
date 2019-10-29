package execution

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type CodeResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type Executor interface {
	Execute(code string) (*CodeResult, error)
}

type baseExecutor struct {
	FileName           string
	CompileCommandName string
	CompileCommandArgs []string
}

func newBaseExecutor(commandName, fileName string) *baseExecutor {
	result := &baseExecutor{
		FileName:           fileName,
		CompileCommandName: commandName,
	}
	result.generateCompileCommandArgs()
	return result
}
func (e *baseExecutor) newFile(code string) error {
	return ioutil.WriteFile(e.FileName, []byte(code), 0644)
}

func (e *baseExecutor) runCommand(name string, arg ...string) (bytes.Buffer, bytes.Buffer) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	return out, errOut
}

func (e *baseExecutor) generateCompileCommandArgs() {
	e.CompileCommandArgs = []string{
		e.FileName,
	}
}

// compile return type is a little bit odd, middle argument is representing
// should you continue because there is no error returned in case of failed compiling
// but you shouldn't continue
func (e *baseExecutor) compile(code string) (*CodeResult, bool, error) {
	err := e.newFile(code)
	if err != nil {
		return nil, false, err
	}
	out, errOut := e.runCommand(e.CompileCommandName, e.CompileCommandArgs...)
	if errOut.Len() != 0 {
		err = os.Remove(e.FileName)
		if err != nil {
			return nil, false, err
		}
		return &CodeResult{
			Output: out.String(),
			Error:  errOut.String(),
		}, false, nil
	}
	return &CodeResult{
		Output: out.String(),
		Error:  errOut.String(),
	}, true, nil
}

type PythonExecutor struct {
	*baseExecutor
}

func NewPythonExecutor() *PythonExecutor {
	return &PythonExecutor{
		baseExecutor: newBaseExecutor("python3", "main.py"),
	}
}

func (e *PythonExecutor) Execute(code string) (*CodeResult, error) {
	result, goOn, err := e.compile(code)
	if !goOn {
		return result, err
	}
	err = os.Remove(e.FileName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ExecPython(text string) (*CodeResult, error) {
	fileName := "main.py"
	err := ioutil.WriteFile(fileName, []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("python3", fileName)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	err = os.Remove(fileName)
	if err != nil {
		return nil, err
	}
	resp := &CodeResult{
		Output: out.String(),
		Error:  errOut.String(),
	}
	return resp, nil
}

type CppExecutor struct {
	*baseExecutor
	RunCommandName string
	ExecutableName string
}

func NewCppExecutor() *CppExecutor {
	result := &CppExecutor{
		baseExecutor:   newBaseExecutor("c++", "main.cpp"),
		ExecutableName: "main",
	}
	result.generateCompileCommandArgs()
	result.generateRunCommandName()
	return result
}

func (e *CppExecutor) generateCompileCommandArgs() {
	e.CompileCommandArgs = []string{
		"-o",
		e.ExecutableName,
		e.FileName,
	}
}

func (e *CppExecutor) generateRunCommandName() {
	e.RunCommandName = fmt.Sprintf("./%v", e.ExecutableName)
}

func (e *CppExecutor) Execute(code string) (*CodeResult, error) {
	result, goOn, err := e.compile(code)
	if !goOn {
		return result, err
	}
	out, errOut := e.runCommand(e.RunCommandName)
	err = os.Remove(e.FileName)
	if err != nil {
		return nil, err
	}
	err = os.Remove(e.ExecutableName)
	if err != nil {
		return nil, err
	}
	return &CodeResult{
		Output: out.String(),
		Error:  errOut.String(),
	}, nil
}

func ExecCpp(text string) (*CodeResult, error) {
	fileName := "main.cpp"
	err := ioutil.WriteFile(fileName, []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	execName := "main"
	cmd := exec.Command("c++", "-o", execName, fileName)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	if errOut.Len() != 0 {
		err = os.Remove(fileName)
		if err != nil {
			return nil, err
		}
		return &CodeResult{
			Error: errOut.String(),
		}, nil
	}
	cmd = exec.Command(fmt.Sprintf("./%v", execName))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	err = os.Remove(fileName)
	err = os.Remove(execName)

	return &CodeResult{
		Output: out.String(),
	}, nil
}

type JavaExecutor struct {
	*baseExecutor
	RunCommandName string
}

func NewJavaExecutor() *JavaExecutor {
	return &JavaExecutor{
		baseExecutor:   newBaseExecutor("javac", "main.java"),
		RunCommandName: "java",
	}
}

func (e *JavaExecutor) Execute(code string) (*CodeResult, error) {
	result, goOn, err := e.compile(code)
	if !goOn {
		return result, err
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}

	unparsedExecName := ""
	for _, file := range files {
		if strings.Contains(file.Name(), ".class") {
			unparsedExecName = file.Name()
			break
		}
	}
	execName := strings.Split(unparsedExecName, ".")[0]

	out, errOut := e.runCommand(e.RunCommandName, execName)
	err = os.Remove(e.FileName)
	if err != nil {
		return nil, err
	}
	err = os.Remove(unparsedExecName)
	if err != nil {
		return nil, err
	}
	return &CodeResult{
		Output: out.String(),
		Error:  errOut.String(),
	}, nil
}

func ExecJava(text string) (*CodeResult, error) {
	fileName := "main.java"
	err := ioutil.WriteFile(fileName, []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("javac", fileName)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	if errOut.Len() != 0 {
		err = os.Remove(fileName)
		if err != nil {
			return nil, err
		}
		return &CodeResult{
			Error: errOut.String(),
		}, nil
	}
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}

	unparsedExecName := ""
	for _, file := range files {
		if strings.Contains(file.Name(), ".class") {
			unparsedExecName = file.Name()
			break
		}
	}
	execName := strings.Split(unparsedExecName, ".")[0]
	cmd = exec.Command("java", execName)
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	err = os.Remove(fileName)
	err = os.Remove(unparsedExecName)

	return &CodeResult{
		Output: out.String(),
	}, nil
}

const (
	Python = "python"
	Java   = "java"
	Cpp    = "c++"
)
