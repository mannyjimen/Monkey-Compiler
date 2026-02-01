package compiler

import (
	"fmt"
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/code"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
	"github.com/mannyjimen/Monkey-Compiler/object"
	"github.com/mannyjimen/Monkey-Compiler/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "3; 4;",
			expectedConstants: []any{3, 4},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2 - 2",
			expectedConstants: []any{2, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMin),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanCompilation(t *testing.T) {
	tests := []compilerTestCase{
		{input: "true;",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{input: "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(bytecode.Instructions, tt.expectedInstructions)
		if err != nil {
			t.Errorf("testInstructions failed: %s", err)
		}

		err = testConstants(bytecode.Constants, tt.expectedConstants)
		if err != nil {
			t.Errorf("testConstants failed: %s", err)
		}
	}
}

func testInstructions(actual code.Instructions, expected []code.Instructions) error {
	flattened := concatInstructions(expected)

	if len(actual) != len(flattened) {
		return fmt.Errorf("incorrect instruction length, \nexpected:%q\ngot\t:%q",
			flattened, actual)
	}

	for i, instr := range actual {
		if flattened[i] != instr {
			return fmt.Errorf("wrong instruction at %d, \nexpected:\t%q\ngot:\t\t%q",
				i, flattened, actual)
		}
	}

	return nil
}

func testConstants(actual []object.Object, expected []any) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("incorrect number of constants, expected %d, got %d",
			len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant at index %d, testIntegerObject failed: %s",
					i, err)
			}
		}
	}

	return nil
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

func concatInstructions(allInstructions []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, instr := range allInstructions {
		out = append(out, instr...)
	}

	return out
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
