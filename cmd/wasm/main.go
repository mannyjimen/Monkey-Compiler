//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/mannyjimen/Monkey-Compiler/compiler"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
	"github.com/mannyjimen/Monkey-Compiler/parser"
	"github.com/mannyjimen/Monkey-Compiler/vm"
)

func runCompiler(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return "error, no input read"
	}

	sourceCode := args[0].String()
	output := runCode(sourceCode)

	return output
}

func runCode(sourceCode string) string {
	var out bytes.Buffer

	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		flattened := strings.Join(p.Errors(), "\n")
		return flattened
	}

	c := compiler.New()
	err := c.Compile(program)
	if err != nil {
		return fmt.Sprintf("compilation failed:\n%s\n", err)
	}

	code := c.Bytecode()

	vm := vm.New(code)
	err = vm.Run()
	if err != nil {
		return fmt.Sprintf("executing bytecode failed:\n%s\n", err)
	}

	lastPopped := vm.LastPoppedObject()

	out.WriteString(lastPopped.Inspect())
	out.WriteString("\n")

	return out.String()
}

func main() {
	js.Global().Set("runMonkey", js.FuncOf(runCompiler))

	<-make(chan bool)
}
