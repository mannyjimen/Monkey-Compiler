package compiler

import (
	"fmt"
	"sort"

	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/code"
	"github.com/mannyjimen/Monkey-Compiler/object"
)

type EmittedInstruction struct {
	Position int
	Opcode   code.Opcode
}

type CompilationScope struct {
	instructions    code.Instructions
	lastInstruction EmittedInstruction
	prevInstruction EmittedInstruction
}

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable

	scope      []CompilationScope
	scopeIndex int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:    code.Instructions{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}

	return &Compiler{
		constants:   []object.Object{},
		symbolTable: NewSymbolTable(),

		scope:      []CompilationScope{mainScope},
		scopeIndex: 0,
	}
}

func NewWithState(symbols *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = symbols
	compiler.constants = constants

	return compiler
}

func (c *Compiler) enterScope() {
	c.scopeIndex++
	c.scope = append(c.scope, CompilationScope{
		instructions:    code.Instructions{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	})
}

func (c *Compiler) exitScope() {
	c.scopeIndex--
	c.scope = c.scope[:len(c.scope)-1]
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		//compiling expression (right side of equal sign)
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

		c.emit(code.OpPop)

	case *ast.BlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown prefix operator %s", node.Operator)
		}

	case *ast.InfixExpression:
		//reorder case
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		// emitting with temporary garbage value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastInstruction()
		}

		//emitting with temp garbage value
		jumpTruthyPos := c.emit(code.OpJump, 9999)

		//modifying jumpNotTruthy
		postConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, postConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIsPop() {
				c.removeLastInstruction()
			}
		}

		postAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpTruthyPos, postAlternativePos)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("identifier %q not found in symbol table", node.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.StringLiteral:
		strLit := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(strLit))

	case *ast.ArrayLiteral:
		arrLen := len(node.Elements)
		for _, elem := range node.Elements {
			err := c.Compile(elem)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, arrLen)

	case *ast.HashLiteral:
		hashLen := len(node.Pairs) * 2

		//sort keys for consistency for tests
		keys := []ast.Expression{}

		for k := range node.Pairs {
			keys = append(keys, k)
		}

		//sort by string values
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, hashLen)

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)

	case *ast.FunctionLiteral:
		funcLit := &object.CompiledFunction{}

		c.enterScope()

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		funcLit.Instructions = c.currentInstructions()

		c.exitScope()

		c.emit(code.OpConstant, c.addConstant(funcLit))

	default:
		return fmt.Errorf("node type not handled: %T", node)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// returns start pos in c.instructions
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instr := code.Make(op, operands...)
	pos := c.addInstruction(instr)

	c.setLastInstruction(op, pos)

	return pos
}

// returns start pos in c.instructions
func (c *Compiler) addInstruction(instr code.Instructions) int {
	posNewInstr := len(c.currentInstructions())
	updatedInstr := append(c.currentInstructions(), instr...)

	c.scope[c.scopeIndex].instructions = updatedInstr

	return posNewInstr
}

// returns index in const pool
func (c *Compiler) addConstant(constant object.Object) int {
	newConstIndex := len(c.constants)
	c.constants = append(c.constants, constant)
	return newConstIndex
}

func (c *Compiler) lastInstructionIsPop() bool {
	last := c.scope[c.scopeIndex].lastInstruction
	return last.Opcode == code.OpPop
}

func (c *Compiler) removeLastInstruction() {
	previous := c.scope[c.scopeIndex].prevInstruction
	last := c.scope[c.scopeIndex].lastInstruction

	oldIns := c.currentInstructions()
	newIns := oldIns[:last.Position]

	c.scope[c.scopeIndex].instructions = newIns
	c.scope[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	newPrev := c.scope[c.scopeIndex].lastInstruction
	newLast := EmittedInstruction{Position: pos, Opcode: op}

	c.scope[c.scopeIndex].prevInstruction = newPrev
	c.scope[c.scopeIndex].lastInstruction = newLast
}

// example of NON type safe function. can possibly corrupt bytecode
// by replacing different length instructions
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := range len(newInstruction) {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	ins := c.currentInstructions()

	op := code.Opcode(ins[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scope[c.scopeIndex].instructions
}
