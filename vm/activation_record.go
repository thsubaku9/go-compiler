package vm

import "monkey-c/code"

type ActivationRecord struct {
	cl                 *code.Closure
	instructionPointer int
	basePointer        int
}

func (ar *ActivationRecord) Instructions() code.Instructions {
	return ar.cl.Fn.Instructions
}

func NewRecord(cl *code.Closure, basePointer int) *ActivationRecord {
	return &ActivationRecord{
		cl:                 cl,
		instructionPointer: -1,
		basePointer:        basePointer,
	}
}
