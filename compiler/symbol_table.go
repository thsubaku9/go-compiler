package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "G"
	LocalScope  SymbolScope = "L"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer   *SymbolTable
	store   map[string]Symbol
	numDefs int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol), Outer: outer}
}

func (symt *SymbolTable) Define(val string) Symbol {

	symbol := Symbol{Name: val, Index: symt.numDefs, Scope: GlobalScope}

	if symt.Outer != nil {
		symbol.Scope = LocalScope
	}
	symt.store[val] = symbol
	symt.numDefs++
	return symbol
}

func (symt *SymbolTable) Resolve(val string) (Symbol, bool) {
	obj, ok := symt.store[val]
	if !ok && symt.Outer != nil {
		obj, ok = symt.Outer.Resolve(val)
	}
	return obj, ok
}
