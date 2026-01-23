package vm

import (
	"fmt"
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/compiler"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
	"github.com/mannyjimen/Monkey-Compiler/object"
	"github.com/mannyjimen/Monkey-Compiler/parser"
)

type vmTestCase struct {
	input    string
	expected any
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		c := compiler.New()
		err := c.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(c.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("runtime error: %s", err)
		}

		stackElem := vm.StackTop()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func TestIngeterArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{input: "1", expected: 1},
		{input: "2", expected: 2},
		{input: "1 + 2", expected: 3},
	}

	runVmTests(t, tests)
}

func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testExpectedObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	actualInt, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object not of type integer, got %T (%+v)", actual, actual)
	}

	if actualInt.Value != expected {
		return fmt.Errorf("incorrect integer value, expected %d, got %d", expected, actualInt.Value)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
