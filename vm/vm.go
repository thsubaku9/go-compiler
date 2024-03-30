package vm

import (
	"fmt"
	"monkey-c/code"
	"monkey-c/compiler"
	"monkey-i/object"
)

const StackLim = 3072

type VM struct {
	instructions code.Instructions
	constants    []object.Object
	stack        []object.Object
	stackPointer int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackLim),
		stackPointer: 0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value
			if op == code.OpAdd {
				result := leftValue + rightValue
				vm.push(&object.Integer{Value: result})
			} else if op == code.OpSub {
				result := leftValue - rightValue
				vm.push(&object.Integer{Value: result})
			}
		case code.OpPop:
			vm.pop()
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

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.stackPointer]
}
