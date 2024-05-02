package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "G"
	LocalScope    SymbolScope = "L"
	FreeScope     SymbolScope = "F"
	FunctionScope SymbolScope = "FN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer       *SymbolTable
	FreeSymbols []Symbol
	store       map[string]Symbol
	numDefs     int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol), FreeSymbols: []Symbol{}}
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

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)
	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1, Scope: FreeScope}
	s.store[original.Name] = symbol
	return symbol
}

func (symt *SymbolTable) Resolve(val string) (Symbol, bool) {
	obj, ok := symt.store[val]
	if !ok && symt.Outer != nil {
		obj, ok = symt.Outer.Resolve(val)

		if !ok || obj.Scope == GlobalScope {
			return obj, ok
		}

		free := symt.defineFree(obj)
		return free, true

	}
	return obj, ok
}
