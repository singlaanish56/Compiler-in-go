package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{"a", GlobalScope, 0},
		"b": Symbol{"b", GlobalScope, 1},
		"c": Symbol{"c", LocalScope, 0},
		"d": Symbol{"d", LocalScope, 1},
		"e": Symbol{"e", LocalScope, 0},
		"f": Symbol{"f", LocalScope, 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected array not found,expected=%+v, got=%+v", a, expected["a"])
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected array not found,expected=%+v, got=%+v", b, expected["b"])
	}

	firstLocal := NewEnclosedSymbolTable(global)
	c := firstLocal.Define("c")
	if c != expected["c"] {
		t.Errorf("expected array not found,expected=%+v, got=%+v", c, expected["c"])
	}
	d := firstLocal.Define("d")
	if d != expected["d"] {
		t.Errorf("expected array not found,expected=%+v, got=%+v", d, expected["d"])
	}
	secondLocal := NewEnclosedSymbolTable(firstLocal)
	e := secondLocal.Define("e")
	if e != expected["e"] {
		t.Errorf("expected array not found,expected=%+v, got=%+v", e, expected["e"])
	}
	f := secondLocal.Define("f")
	if f != expected["f"] {
		t.Errorf("expected array not found,expected=%+v, got=%+v", f, expected["f"])
	}
}

func TestResolve(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		Symbol{"a", GlobalScope, 0},
		Symbol{"b", GlobalScope, 1},
	}

	for _, got := range expected {
		result, ok := global.Resolve(got.Name)
		if !ok {
			t.Errorf("name %s is not found in the symbol table", got.Name)
			continue
		}
		if result != got {
			t.Errorf("expected %s to resolve to %+v, got=%+v", got.Name, got, result)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		Symbol{"a", GlobalScope, 0},
		Symbol{"b", GlobalScope, 1},
		Symbol{"c", LocalScope, 0},
		Symbol{"d", LocalScope, 1},
	}

	for _, got := range expected {
		result, ok := local.Resolve(got.Name)
		if !ok {
			t.Errorf("name %s is not found in the symbol table", got.Name)
			continue
		}
		if result != got {
			t.Errorf("expected %s to resolve to %+v, got=%+v", got.Name, got, result)
		}
	}
}

func TestDefineResolveWithBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		Symbol{"a", BuiltinScope, 0},
		Symbol{"c", BuiltinScope, 1},
		Symbol{"e", BuiltinScope, 2},
		Symbol{"f", BuiltinScope, 3},
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, got := range expected {
			result, ok := table.Resolve(got.Name)
			if !ok {
				t.Errorf("name %s is not found in the symbol table", got.Name)
				continue
			}
			if result != got {
				t.Errorf("expected %s to resolve to %+v, got=%+v", got.Name, got, result)
			}

		}
	}
}
