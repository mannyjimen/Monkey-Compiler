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
				code.Make(code.OpSub),
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

func TestBooleanExpressions(t *testing.T) {
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
		{input: "1 > 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{input: "1 < 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{input: "1 == 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{input: "1 != 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestPrefixExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{input: "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
		{input: "-5",
			expectedConstants: []any{5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "if (true) { 1 }; 2;",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 11),
				// 000(10)
				code.Make(code.OpNull),
				// 000(11)
				code.Make(code.OpPop),
				// 000(12)
				code.Make(code.OpConstant, 1),
				// 000(15)
				code.Make(code.OpPop),
			},
		},
		{
			input:             "if (5 > 3) { 0 }",
			expectedConstants: []any{5, 3, 0},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpConstant, 0),
				// 0003
				code.Make(code.OpConstant, 1),
				// 0006
				code.Make(code.OpGreaterThan),
				// 0007
				code.Make(code.OpJumpNotTruthy, 16),
				// 000A (10)
				code.Make(code.OpConstant, 2),
				// 000D (13)
				code.Make(code.OpJump, 17),
				//000(16)
				code.Make(code.OpNull),
				//000(17)
				code.Make(code.OpPop),
			},
		},
		{
			input:             "if (true) { 0 } else { 1 }",
			expectedConstants: []any{0, 1},
			expectedInstructions: []code.Instructions{
				//0000
				code.Make(code.OpTrue),
				//0001
				code.Make(code.OpJumpNotTruthy, 10),
				//0004
				code.Make(code.OpConstant, 0),
				//0007
				code.Make(code.OpJump, 13),
				//000(10)
				code.Make(code.OpConstant, 1),
				//000(13)
				code.Make(code.OpPop),
			},
		},

		{
			input:             "if (true) { 0;1;2; } else { 3;4;5 }",
			expectedConstants: []any{0, 1, 2, 3, 4, 5},
			expectedInstructions: []code.Instructions{
				//0000
				code.Make(code.OpTrue),
				//0001
				code.Make(code.OpJumpNotTruthy, 18),
				//0004
				code.Make(code.OpConstant, 0),
				//0007
				code.Make(code.OpPop),
				//0008
				code.Make(code.OpConstant, 1),
				//000(11)
				code.Make(code.OpPop),
				//000(12)
				code.Make(code.OpConstant, 2),
				//000(15)
				code.Make(code.OpJump, 29),
				//000(18)
				code.Make(code.OpConstant, 3),
				//000(21)
				code.Make(code.OpPop),
				//000(22)
				code.Make(code.OpConstant, 4),
				//000(25)
				code.Make(code.OpPop),
				//000(26)
				code.Make(code.OpConstant, 5),
				//000(29)
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "let x = 5;",
			expectedConstants: []any{5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				//x is 0
				code.Make(code.OpSetGlobal, 0),
			},
		},
		{
			input: `
			let first = 5;
			let second = first;`,
			expectedConstants: []any{5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				//first is 0
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				//second is 1
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
			let x = true;
			let y = x;
			x; y;`,
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				//x is 0
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				//y is 1
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			`"hello"`,
			[]any{"hello"},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			`"mon" + "key"`,
			[]any{"mon", "key"},
			[]code.Instructions{
				code.Make(code.OpConstant, 0), //"mon"
				code.Make(code.OpConstant, 1), //"key"
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			`let x = "hello"`,
			[]any{"hello"},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),  //hello
				code.Make(code.OpSetGlobal, 0), // let x =
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			`[]`,
			[]any{},
			[]code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			`["hello"]`,
			[]any{"hello"},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),

				code.Make(code.OpArray, 1),
				code.Make(code.OpPop),
			},
		},
		{
			`[1 + 2, "mon" + "key"]`,
			[]any{1, 2, "mon", "key"},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpAdd),

				code.Make(code.OpArray, 2),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			`{}`,
			[]any{},
			[]code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			`{1:2}`,
			[]any{1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0), //1
				code.Make(code.OpConstant, 1), //2

				code.Make(code.OpHash, 2),
				code.Make(code.OpPop),
			},
		},
		{
			`{2 + 2 : 4 * 1, 1: 2}`,
			[]any{2, 2, 4, 1, 1, 2},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpMul),

				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),

				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexOperatorExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			`[1, 2, 3][1]`,
			[]any{1, 2, 3, 1},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),

				code.Make(code.OpArray, 3), //array

				code.Make(code.OpConstant, 3), //index

				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			`[1 + 1, 2, 5 * 5][1 + 1]`,
			[]any{1, 1, 2, 5, 5, 1, 1},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpMul),

				code.Make(code.OpArray, 3), //array

				code.Make(code.OpConstant, 5),
				code.Make(code.OpConstant, 6),
				code.Make(code.OpAdd),

				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			`{1 + 1: 2, 3: 5 * 5}[1 + 1]`,
			[]any{1, 1, 2, 3, 5, 5, 1, 1},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),

				code.Make(code.OpHash, 4), //hash

				code.Make(code.OpConstant, 6),
				code.Make(code.OpConstant, 7),
				code.Make(code.OpAdd),

				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			`fn() { return 5 + 10; }`,
			[]any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
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
				return fmt.Errorf("constant at index %d - testIntegerObject failed: %s",
					i, err)
			}

		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant at index %d - testStringObject failed: %s",
					i, err)
			}
		case []code.Instructions:
			compiledFn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant at index %d, expected COMPILED_FUNCTION, got %s", i, actual[i].Type())
			}

			err := testInstructions(compiledFn.Instructions, constant)
			if err != nil {
				return fmt.Errorf("constant at index %d - testInstructions failed: %s", i, err)
			}

		default:
			return fmt.Errorf("constant type at index %d not handled, type: %T",
				i, constant)
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	actualInt, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("compiled object not of type integer, got %T (%+v)", actual, actual)
	}

	if actualInt.Value != expected {
		return fmt.Errorf("incorrect integer value, expected %d, got %d", expected, actualInt.Value)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	actualStr, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("compiled object not of type string, got %T (%+v)", actual, actual)
	}

	if actualStr.Value != expected {
		return fmt.Errorf("incorrect string value, expected %q, got %q",
			expected, actualStr.Value)
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
