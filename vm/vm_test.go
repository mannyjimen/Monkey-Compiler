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
			t.Fatalf("%s", err)
		}

		vm := New(c.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("%s", err)
		}
		stackElem := vm.LastPoppedObject()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func TestIngeterArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{input: "1", expected: 1},
		{input: "2", expected: 2},
		{input: "1 + 2", expected: 3},
		{input: "5 - 5", expected: 0},
		{input: "1 * 2", expected: 2},
		{input: "2 / 1", expected: 2},
		{input: "50 + 6 - 4 * 3 * 2 - 2 / 2", expected: 31},
		{input: "9 * (4 - 2)", expected: 18},
		{input: "-1", expected: -1},
		{input: "-(9 * (4 - 2))", expected: -18},
	}

	runVmTests(t, tests)
}

func TestBooleanInstructions(t *testing.T) {
	tests := []vmTestCase{
		{input: "true", expected: true},
		{input: "false;", expected: false},
		{input: "1 > 2", expected: false},
		{input: "1 < 2", expected: true},
		{input: "1 == 2", expected: false},
		{input: "1 != 2", expected: true},
		{input: "true == false", expected: false},
		{input: "true == true", expected: true},
		{input: "(true != true)", expected: false},
		{input: "(false == false)", expected: true},
		{input: "!true", expected: false},
		{input: "!(!(true))", expected: true},
		{input: "!!!false", expected: true},
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
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	actualInt, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object not of type Integer, got %T (%+v)", actual, actual)
	}

	if actualInt.Value != expected {
		return fmt.Errorf("incorrect integer value, expected %d, got %d", expected, actualInt.Value)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	actualBool, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object not of type Boolean, got %T (%+v)", actual, actual)
	}

	if actualBool.Value != expected {
		return fmt.Errorf("incorrect boolean value, expected %t, got %t", expected, actualBool.Value)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
