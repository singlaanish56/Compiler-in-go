package compiler

type SymbolScope string

const (
	GloablScope SymbolScope = "GLOBAL"
)

type Symbol struct{
	Name string
	Scope SymbolScope
	Position int
}

type SymbolTable struct{
	store map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable{
	return &SymbolTable{
		store: make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func (st *SymbolTable) Define(name string) Symbol{
	symbol := Symbol{name, GloablScope, st.numDefinitions}
	st.store[name]=symbol
	st.numDefinitions++
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool){
	obj, ok := st.store[name]
	return obj, ok
}	