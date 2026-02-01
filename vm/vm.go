package vm

import (
	"fmt"

	"github.com/mannyjimen/Monkey-Compiler/code"
	"github.com/mannyjimen/Monkey-Compiler/compiler"
	"github.com/mannyjimen/Monkey-Compiler/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int //points to next free stack position, sp - 1 is top of stack index
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// hot-path
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return fmt.Errorf("runtime error: %s\n", err)
			}
		case code.OpAdd, code.OpMin, code.OpMul, code.OpDiv:

			err := vm.executeInfixOperation(op)
			if err != nil {
				return fmt.Errorf("runtime error: %s\n", err)
			}

		case code.OpPop:
			_, err := vm.pop()
			if err != nil {
				return fmt.Errorf("runtime error: %s\n", err)
			}
		}
	}

	return nil
}

func (vm *VM) executeInfixOperation(operator code.Opcode) error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	leftInteger, leftValid := left.(*object.Integer)
	rightInteger, rightValid := right.(*object.Integer)

	if !leftValid || !rightValid {
		return fmt.Errorf("left or right operand for infix operation is not Integer\n")
	}

	var result int64

	switch operator {
	case code.OpAdd:
		result = leftInteger.Value + rightInteger.Value
	case code.OpMin:
		result = leftInteger.Value - rightInteger.Value
	case code.OpMul:
		result = leftInteger.Value * rightInteger.Value
	case code.OpDiv:
		result = leftInteger.Value / rightInteger.Value
	default:
		return fmt.Errorf("unknown infix operator: %d", operator)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

// note, remove conditional if speed is prioritized
func (vm *VM) pop() (object.Object, error) {
	if vm.sp <= 0 {
		return nil, fmt.Errorf("stack underflow")
	}

	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj, nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// test only method
func (vm *VM) LastPoppedObject() object.Object {
	return vm.stack[vm.sp]
}
