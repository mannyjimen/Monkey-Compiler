package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/mannyjimen/Monkey-Compiler/compiler"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
	"github.com/mannyjimen/Monkey-Compiler/parser"
	"github.com/mannyjimen/Monkey-Compiler/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		c := compiler.New()
		err := c.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "compilation failed:\n%s\n", err)
			continue
		}

		vm := vm.New(c.Bytecode())
		err = vm.Run()
		if err != nil {
			fmt.Fprintf(out, "executing bytecode failed:\n%s\n", err)
			continue
		}

		io.WriteString(out, vm.StackTop().Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
