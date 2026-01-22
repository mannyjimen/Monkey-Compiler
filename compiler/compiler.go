package compiler

import (
	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/code"
	"github.com/mannyjimen/Monkey-Compiler/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
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
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		//emit OpAdd instr should go here in the future
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
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
