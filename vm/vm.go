package vm

import (
	"fmt"
	"monkey-c/code"
	"monkey-c/compiler"
	"monkey-i/object"
)

const (
	StackLim    = 2048
	GlobalsSize = 2048
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	instructions code.Instructions
	constants    []object.Object
	stack        []object.Object
	stackPointer int
	globals      []object.Object
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackLim),
		stackPointer: 0,
		globals:      make([]object.Object, GlobalsSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			if err := vm.push(vm.constants[constIndex]); err != nil {
				return err
			}

		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.executeBinaryOperation(op); err != nil {
				return err
			}

		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}

		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}

		case code.OpBang:
			if err := vm.executeBangOperator(); err != nil {
				return err
			}

		case code.OpMinus:
			if err := vm.executeNumberNegation(); err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanEqual:
			if err := vm.executeComp(op); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.OpSetGlobal:
			gIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[gIdx] = vm.pop()

		case code.OpGetGlobal:
			gIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			if err := vm.push(vm.globals[gIdx]); err != nil {
				return err
			}

		case code.OpArray:
			numElems := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			arr := vm.buildArray(vm.stackPointer-numElems, vm.stackPointer)
			if err := vm.push(arr); err != nil {
				return err
			}
		}

	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.stackPointer >= StackLim {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.stackPointer] = o
	vm.stackPointer++
	return nil
}

func (vm *VM) pop() object.Object {
	if vm.stackPointer == 0 {
		panic("stack is empty !!")
	}

	o := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	return o
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.stackPointer]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {

	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left.(*object.Integer), right.(*object.Integer))
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left.(*object.String), right.(*object.String))
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right *object.Integer) error {
	var result int64
	switch op {
	case code.OpAdd:
		result = left.Value + right.Value
	case code.OpSub:
		result = left.Value - right.Value
	case code.OpMul:
		result = left.Value * right.Value
	case code.OpDiv:
		result = left.Value / right.Value
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right *object.String) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	return vm.push(&object.String{Value: left.Value + right.Value})
}

func (vm *VM) executeComp(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left.(*object.Integer), right.(*object.Integer))
	} else if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return vm.push(nativeBoolToBooleanObject(left.(*object.String).Value == right.(*object.String).Value))

	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right *object.Integer) error {
	leftValue := left.Value
	rightValue := right.Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	case code.OpGreaterThanEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue >= rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (vm *VM) executeBangOperator() error {

	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False, Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeNumberNegation() error {

	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elems := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elems[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elems}
}
