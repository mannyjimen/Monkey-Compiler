package compiler

import (
	"fmt"

	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/code"
	"github.com/mannyjimen/Monkey-Compiler/object"
)

type EmittedInstruction struct {
	Position int
	Opcode   code.Opcode
}

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction EmittedInstruction
	prevInstruction EmittedInstruction
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},

		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}
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

		//what now?

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
		postConsequencePos := len(c.instructions)
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

		postAlternativePos := len(c.instructions)
		c.changeOperand(jumpTruthyPos, postAlternativePos)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	default:
		return fmt.Errorf("node type not handled: %T", node)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
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
	posNewInstr := len(c.instructions)
	c.instructions = append(c.instructions, instr...)
	return posNewInstr
}

// returns index in const pool
func (c *Compiler) addConstant(constant object.Object) int {
	c.constants = append(c.constants, constant)
	return len(c.constants) - 1
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastInstruction() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.prevInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.prevInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{Position: pos, Opcode: op}
}

// example of NON type safe function. can possibly corrupt bytecode
// by replacing different length instructions
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := range len(newInstruction) {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}
