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
	result, compiled, err := e.compile(code)
	if !compiled {
		return result, err
	}
	err = os.Remove(e.FileName)
	if err != nil {
		return nil, err
	}
	return result, nil
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
	result, compiled, err := e.compile(code)
	if !compiled {
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

func (e *JavaExecutor) classFile() string {
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
	return unparsedExecName
}

func (e *JavaExecutor) Execute(code string) (*CodeResult, error) {
	result, compiled, err := e.compile(code)
	if !compiled {
		return result, err
	}

	unparsedExecName := e.classFile()
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

const (
	Python = "python"
	Java   = "java"
	Cpp    = "c++"
)
