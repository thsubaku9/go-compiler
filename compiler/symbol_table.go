package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "G"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store   map[string]Symbol
	numDefs int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func (symt *SymbolTable) Define(val string) Symbol {

	symbol := Symbol{Name: val, Index: symt.numDefs, Scope: GlobalScope}
	symt.store[val] = symbol
	symt.numDefs++
	return symbol
}

func (symt *SymbolTable) Resolve(val string) (Symbol, bool) {
	obj, ok := symt.store[val]
	return obj, ok
}
