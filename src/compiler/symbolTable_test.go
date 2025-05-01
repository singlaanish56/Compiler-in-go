package compiler

import "testing"

func TestDefine(t *testing.T){
	 expected := map[string]Symbol{
		"a":Symbol{"a", GloablScope, 0},
		"b":Symbol{"b", GloablScope, 1},
	 }

	 global := NewSymbolTable()

	 a := global.Define("a")
	 if a != expected["a"]{
		t.Errorf("expected array not found,expected=%+v, got=%+v",a, expected["a"])
	 }

	 b := global.Define("b")
	 if b != expected["b"]{
		t.Errorf("expected array not found,expected=%+v, got=%+v",b, expected["b"])
	 }
}

func TestResolve(t *testing.T){
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		Symbol{"a", GloablScope, 0},
		Symbol{"b", GloablScope, 1},
	}

	for _, got := range expected{
		result, ok := global.Resolve(got.Name)
		if !ok{
			t.Errorf("name %s is not found in the symbol table", got.Name)
			continue
		}
		if result != got{
			t.Errorf("expected %s to resolve to %+v, got=%+v", got.Name, got, result)
		}
	}
}