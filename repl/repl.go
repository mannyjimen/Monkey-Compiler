package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/mannyjimen/Monkey-Compiler/compiler"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
	"github.com/mannyjimen/Monkey-Compiler/object"
	"github.com/mannyjimen/Monkey-Compiler/parser"
	"github.com/mannyjimen/Monkey-Compiler/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	//slices, we hold on to these to remember STATE
	constants := []object.Object{}
	symbolTable := compiler.NewSymbolTable()
	globals := make([]object.Object, vm.GlobalsSize)

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		// fmt.Println(scanned)
		if !scanned || len(scanner.Text()) == 0 { //no more user input
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		c := compiler.NewWithState(symbolTable, constants)
		err := c.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "compilation failed:\n%s\n", err)
			continue
		}

		code := c.Bytecode()
		constants = code.Constants

		vm := vm.NewWithGlobalsStore(code, globals)
		err = vm.Run()
		if err != nil {
			fmt.Fprintf(out, "executing bytecode failed:\n%s\n", err)
			continue
		}

		lastPopped := vm.LastPoppedObject()

		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
