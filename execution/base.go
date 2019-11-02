package execution

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Languages constants
const (
	Python = "python"
	Java   = "java"
	Cpp    = "cpp"
)

// Main file names
const (
	PythonMainName = "main.py"
	JavaMainName   = "main.java"
	CppMainName    = "main.cpp"
)

// CodeResult have output and error
type CodeResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

// Executor executes code
type Executor interface {
	Execute() (*CodeResult, error)
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

func hasError(errString string) bool {
	lowerErrString := strings.ToLower(errString)
	if strings.Contains(lowerErrString, "error") {
		return true
	}
	if strings.Contains(lowerErrString, "traceback") {
		return true
	}
	return false
}

// compile return type is a little bit odd, middle argument is representing
// should you continue because there is no error returned in case of failed compiling
// but you shouldn't continue
func (e *baseExecutor) compile() (*CodeResult, bool, error) {
	out, errOut := e.runCommand(e.CompileCommandName, e.CompileCommandArgs...)
	if hasError(errOut.String()) {
		if err := os.Remove(e.FileName); err != nil {
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

// PythonExecutor is Executor for python code
type PythonExecutor struct {
	*baseExecutor
}

// NewPythonExecutor initializes new instance of PythonExecutor
func NewPythonExecutor() *PythonExecutor {
	return &PythonExecutor{
		baseExecutor: newBaseExecutor("python3", PythonMainName),
	}
}

// Execute executes code
func (e *PythonExecutor) Execute() (*CodeResult, error) {
	result, compiled, err := e.compile()
	if !compiled {
		return result, err
	}
	if err = os.Remove(e.FileName); err != nil {
		return nil, err
	}
	return result, nil
}

// CppExecutor is Executor for c++ code
type CppExecutor struct {
	*baseExecutor
	RunCommandName string
	ExecutableName string
}

// NewCppExecutor initializes new instance of CppExecutor
func NewCppExecutor() *CppExecutor {
	result := &CppExecutor{
		baseExecutor:   newBaseExecutor("c++", CppMainName),
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

// Execute executes code
func (e *CppExecutor) Execute() (*CodeResult, error) {
	result, compiled, err := e.compile()
	if !compiled {
		return result, err
	}
	out, errOut := e.runCommand(e.RunCommandName)
	if err = os.Remove(e.FileName); err != nil {
		return nil, err
	}
	if err = os.Remove(e.ExecutableName); err != nil {
		return nil, err
	}
	return &CodeResult{
		Output: out.String(),
		Error:  errOut.String(),
	}, nil
}

// JavaExecutor is Executor for java code
type JavaExecutor struct {
	*baseExecutor
	RunCommandName string
}

// NewJavaExecutor initializes new instance of JavaExecutor
func NewJavaExecutor() *JavaExecutor {
	return &JavaExecutor{
		baseExecutor:   newBaseExecutor("javac", JavaMainName),
		RunCommandName: "java",
	}
}

func (e *JavaExecutor) classFiles() []string {
	files, _ := ioutil.ReadDir("./")
	var fileNames []string
	for _, file := range files {
		if strings.Contains(file.Name(), ".class") {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func (e *JavaExecutor) removeClassFiles() {
	classNames := e.classFiles()
	for _, className := range classNames {
		os.Remove(className)
	}
}

// Execute executes code
func (e *JavaExecutor) Execute() (*CodeResult, error) {
	result, compiled, err := e.compile()
	if !compiled {
		e.removeClassFiles()
		return result, err
	}
	out, errOut := e.runCommand(e.RunCommandName, "Main")
	if err = os.Remove(e.FileName); err != nil {
		return nil, err
	}
	e.removeClassFiles()
	return &CodeResult{
		Output: out.String(),
		Error:  errOut.String(),
	}, nil
}
