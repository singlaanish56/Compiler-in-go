package vm

import (
	"fmt"
	"testing"
	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/compiler"
	"github.com/singlaanish56/Compiler-in-go/lexer"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/parser"
)


func parse(input string) *ast.AstRootNode{
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParserProgram()
}

func testIntegerObject(expected int64, obj object.Object) error{
	result, ok := obj.(*object.Integer)
	if !ok{
		return fmt.Errorf("object is not an Integer, got=%T(%+v)", obj, obj)
	}

	if result.Value != expected{
		return fmt.Errorf("object has wrong value, got=%d, want=%d", result.Value, expected)
	}

	return nil
}

type vmTestCase struct{
	input string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase){
	t.Helper()

	for _, tt := range tests{
		 program := parse(tt.input)

		 compiler := compiler.New()
		 err := compiler.Compile(program)
		 if err != nil{
			t.Fatalf("failed to compile: %s", err)
		}

		vm := New(compiler.Bytecode())
		err = vm.Run()
		if err != nil{
			t.Fatalf("failed to run vm: %s", err)
		}

		stackTop := vm.StackTop()

		testExpectedObject(t, tt.expected, stackTop)

	}
}

func testExpectedObject(t *testing.T, expected interface{}, obj object.Object){
	t.Helper()

	switch expected := expected.(type){
	case int:
		err := testIntegerObject(int64(expected), obj)
		if err != nil{
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T){
	tests := []vmTestCase{
		{"1",1},
		{"2",2},
		{"1+2",3},
	}

	runVmTests(t, tests)
}