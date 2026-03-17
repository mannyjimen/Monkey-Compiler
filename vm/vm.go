package vm

import (
	"fmt"

	"github.com/mannyjimen/Monkey-Compiler/code"
	"github.com/mannyjimen/Monkey-Compiler/compiler"
	"github.com/mannyjimen/Monkey-Compiler/object"
)

// caps
const StackSize = 2048
const GlobalsSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int //points to next free stack position, sp - 1 is top of stack index

	globals []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
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
				return formatRuntimeError(err)
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:

			err := vm.executeInfixOperation(op)
			if err != nil {
				return formatRuntimeError(err)
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return formatRuntimeError(err)
			}

		case code.OpBang, code.OpMinus:
			err := vm.executePrefix(op)
			if err != nil {
				return formatRuntimeError(err)
			}

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return formatRuntimeError(err)
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return formatRuntimeError(err)
			}

		case code.OpJump:
			jumpIndex := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = jumpIndex - 1

		case code.OpJumpNotTruthy:
			jumpIndex := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition, err := vm.pop()
			if err != nil {
				return formatRuntimeError(err)
			}

			if !isTruthy(condition) {
				ip = jumpIndex - 1
			}

		case code.OpPop:
			_, err := vm.pop()
			if err != nil {
				return formatRuntimeError(err)
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return formatRuntimeError(err)
			}

		case code.OpSetGlobal:
			globalsIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			obj, err := vm.pop()
			if err != nil {
				return formatRuntimeError(err)
			}

			vm.globals[globalsIndex] = obj

		case code.OpGetGlobal:
			globalsIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.globals[globalsIndex])
			if err != nil {
				return formatRuntimeError(err)
			}

		default:
			return formatRuntimeError(fmt.Errorf("unknown opcode %d", op))
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

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeInfixIntegerOperation(operator, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeInfixStringOperation(operator, left, right)
	default:
		return fmt.Errorf("unsupported or mismatched infix operand types: %s and %s",
			leftType, rightType)
	}

}

func (vm *VM) executeInfixIntegerOperation(operator code.Opcode, left object.Object, right object.Object) error {
	var result int64

	leftInteger := left.(*object.Integer)
	rightInteger := right.(*object.Integer)

	switch operator {
	case code.OpAdd:
		result = leftInteger.Value + rightInteger.Value
	case code.OpSub:
		result = leftInteger.Value - rightInteger.Value
	case code.OpMul:
		result = leftInteger.Value * rightInteger.Value
	case code.OpDiv:
		result = leftInteger.Value / rightInteger.Value
	default:
		return fmt.Errorf("unknown infix operator %d for operands of type INTEGER", operator)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeInfixStringOperation(operator code.Opcode, left object.Object, right object.Object) error {
	var result string

	leftStr := left.(*object.String)
	rightStr := right.(*object.String)

	switch operator {
	case code.OpAdd:
		result = leftStr.Value + rightStr.Value
	default:
		return fmt.Errorf("unsupported infix operator %d for operands of type STRING", operator)
	}

	return vm.push(&object.String{Value: result})
}

func (vm *VM) executeComparison(operator code.Opcode) error {
	right, err := vm.pop()
	if err != nil {
		return err
	}

	left, err := vm.pop()
	if err != nil {
		return err
	}

	if left.Type() == object.INTEGER_OBJ || right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(operator, left, right)
	}

	//we now know both left and right are our boolean objects (True and False)

	switch operator {
	case code.OpEqual:
		return vm.push(nativeBooleanToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBooleanToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator type: %T, operand types: (%s, %s)",
			operator, left.Type(), right.Type())
	}
}

func (vm *VM) executePrefix(operator code.Opcode) error {
	value, err := vm.pop()
	if err != nil {
		return err
	}

	switch operator {
	case code.OpBang:
		return vm.executeBangOperator(value)

	case code.OpMinus:
		i, ok := value.(*object.Integer)
		if !ok {
			return fmt.Errorf("'-' prefix operator not followed by INTEGER, followed by %s", value.Type())
		}
		return vm.push(&object.Integer{Value: -i.Value})

	default:
		return fmt.Errorf("unknown prefix operator: %T", operator)
	}
}

func (vm *VM) executeBangOperator(value object.Object) error {
	switch value := value.(type) {
	case *object.Null:
		return vm.push(True)
	case *object.Boolean:
		return vm.push(nativeBooleanToBooleanObject(!value.Value))
	default:
		return fmt.Errorf("unexpected type after '!' operator, type %s", value.Type())
	}
}

func (vm *VM) executeIntegerComparison(operator code.Opcode, left, right object.Object) error {
	leftInteger, leftValid := left.(*object.Integer)
	rightInteger, rightValid := right.(*object.Integer)

	if !leftValid || !rightValid {
		return fmt.Errorf("unexpected type for integer comparison, left type: %s, right type: %s",
			left.Type(), right.Type())
	}

	switch operator {
	case code.OpEqual:
		return vm.push(nativeBooleanToBooleanObject(leftInteger.Value == rightInteger.Value))
	case code.OpNotEqual:
		return vm.push(nativeBooleanToBooleanObject(leftInteger.Value != rightInteger.Value))
	case code.OpGreaterThan:
		return vm.push(nativeBooleanToBooleanObject(leftInteger.Value > rightInteger.Value))
	default:
		return fmt.Errorf("unknown operator type: %T, operand types: (%s, %s)",
			operator, left.Type(), right.Type())
	}
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
		// return Null, nil
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

func nativeBooleanToBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	//0 is not truthy to me!
	case *object.Integer:
		if obj.Value == 0 {
			return false
		}
		return true
	default:
		return true
	}
}

func formatRuntimeError(err error) error {
	return fmt.Errorf("runtime error: %s\n", err)
}
