package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name     string
	Scope    SymbolScope
	Position int
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (st *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{name, GlobalScope, st.numDefinitions}
	if st.Outer ==nil{
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	st.store[name] = symbol
	st.numDefinitions++
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := st.store[name]
	if !ok && st.Outer!=nil{
		obj, ok = st.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}
