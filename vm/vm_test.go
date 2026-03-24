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
		{input: "!(if (false) { 1; })", expected: true},
	}

	runVmTests(t, tests)
}

func TestIfExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: "if (true) { 1 }", expected: 1},
		{input: "if (false) { 1 } else { 2 }", expected: 2},
		{input: "if (true) { 1 } else { 2 }", expected: 1},
		{input: "if (1 < 2) { 5 } else { 3 }", expected: 5},
		{input: "if (1 > 2) { 5 } else { 3 }", expected: 3},
		{input: "if (1) { 2 }", expected: 2},
		{input: "if (false) { 2 }", expected: Null},
		{input: "if (1 == 2) { 1 };", expected: Null},
		{input: "!(if (if (1 < 2) { false } else { true }) { 5 })", expected: true},
		{input: "if (if (2 < 1) { true }) { 5 } else { 4 }", expected: 4},
		{input: "if (0) { true } else { false }", expected: false},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{input: "let one = 1; one;", expected: 1},
		{input: "let one = 1; let two = 2; one + two;", expected: 3},
		{input: "let one = 1; let two = one + one; one + two;", expected: 3},
		{input: "let four = 4; let four2 = four; (four + four2) / four2 ;", expected: 2},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: `"hello"`, expected: "hello"},
		{input: `"foo" + "bar"`, expected: "foobar"},
		{input: `"do" + "re" + "mi"`, expected: "doremi"},
	}

	runVmTests(t, tests)
}

func TestArrayExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: `[]`, expected: []int{}},
		{input: `[1, 2, 3]`, expected: []int{1, 2, 3}},
		{input: `[1 + 2, 3 + 4]`, expected: []int{3, 7}},
	}

	runVmTests(t, tests)
}

func TestHashExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: `{}`, expected: map[object.HashKey]int64{}},
		{
			input:    `{1 : 2}`,
			expected: map[object.HashKey]int64{(&object.Integer{Value: 1}).HashKey(): 2},
		},
		{
			input: `{1 + 2: 1 * 4, 1 : 2}`,
			expected: map[object.HashKey]int64{
				(&object.Integer{Value: 3}).HashKey(): 4,
				(&object.Integer{Value: 1}).HashKey(): 2},
		},
		{
			input: `{"hello": 1 * 4, 1 : 2}`,
			expected: map[object.HashKey]int64{
				(&object.String{Value: "hello"}).HashKey(): 4,
				(&object.Integer{Value: 1}).HashKey():      2},
		},
	}

	runVmTests(t, tests)
}

func TestIndexOperatorExpressions(t *testing.T) {
	tests := []vmTestCase{
		{input: `[1][0]`, expected: 1},
		{input: `[1, 2 + 2][1]`, expected: 4},
		{input: `{1 : 1, "monkey" : "hello"}["mon" + "key"]`, expected: "hello"},

		//harsher cases
		{input: `[][0]`, expected: Null},
		{input: `[1, 2, 3][50]`, expected: Null},
		{input: `[1, 2, 3][-1]`, expected: Null},
		{input: `{"hey": "goodbye"}["monkey"]`, expected: Null},
		{input: `[[1, 2], 3][0][0]`, expected: 1},
		{input: `{}[0]`, expected: Null},
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
			t.Errorf("testExpectedObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testExpectedObject failed: %s", err)
		}

	case []int:
		arrLit, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("testExpectedObject failed: expected object.Array, got %T", actual)
		}

		if len(arrLit.Elements) != len(expected) {
			t.Errorf("incorrect number of Array elements, expected %d, got %d",
				len(expected), len(arrLit.Elements))
		}

		for i, num := range arrLit.Elements {
			err := testIntegerObject(int64(expected[i]), num)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case map[object.HashKey]int64:
		hashLit, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("testExpectedObject failed: expected object.Hash, got %T", actual)
		}

		if len(hashLit.Pairs) != len(expected) {
			t.Errorf("incorrect number of Hash elements, expected %d, got %d",
				len(expected), len(hashLit.Pairs))
		}

		for expectedKey, expectedVal := range expected {
			hashPair, ok := hashLit.Pairs[expectedKey]
			if !ok {
				t.Errorf("expected key does not exist %+v", expectedKey)
			}

			err := testIntegerObject(expectedVal, hashPair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case *object.Null:
		if actual != Null {
			t.Errorf("testExpecedObject failed: object not of type Null, got %T (%+v)", actual, actual)
		}

	default:
		t.Errorf("testExpectedObject failed: unhandled expected object type: %T", expected)

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

func testStringObject(expected string, actual object.Object) error {
	actualStr, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object not of type String, got %T (%+v)", actual, actual)
	}

	if actualStr.Value != expected {
		return fmt.Errorf("incorrect string value, expected %q, got %q",
			expected, actualStr.Value)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
