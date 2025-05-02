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

func testBooleanObject(expected bool, obj object.Object) error{
	result, ok := obj.(*object.Boolean)
	if !ok{
		return fmt.Errorf("object is not a Boolean, got=%T(%+v)", obj, obj)
	}

	if result.Value != expected{
		return fmt.Errorf("object has wrong value, got=%t, want=%t", result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, obj object.Object) error{
	result, ok := obj.(*object.String)
	if !ok{
		return fmt.Errorf("the object type is not string, got=%T", obj)
	}

	if result.Value != expected{
		return fmt.Errorf("the expected=%s string value is not equal, got=%s", expected, result.Value)
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

		stackTop := vm.LastPoppedStackElement()

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
	case bool:
		err := testBooleanObject(bool(expected), obj)
		if err != nil{
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, obj)
		if err != nil{
			t.Errorf("test for the string object failed: %s", err)
		}
	case *object.Null:
		if obj!=Null{
			t.Errorf("object is not Null, got=%T(%+v)", obj, obj)
		}
	}
}

func TestIntegerArithmetic(t *testing.T){
	tests := []vmTestCase{
		{"1",1},
		{"2",2},
		{"1+2",3},
		{"2-1", 1},
		{"2*3", 6},
		{"6/2", 3},
		{"1*3-4", -1},
		{"5 * (2+10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50+10", -40},
	}

	runVmTests(t, tests)
}

func TestBooleanLiterals(t *testing.T){
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1>2", false},
		{"1<2", true},
		{"5==4", false},
		{"5!=4", true},
		{"true==false", false},
		{"true!=false", true},
		{"!true", false},
		{"!!true", true},
		{"!!false", false},
		{"!false", true},
	}

	runVmTests(t, tests)
}


func TestConditionals(t *testing.T){
	tests := []vmTestCase{
		{"if(true){10}", 10},
		{"if(true){10}else{20}", 10},
		{"if(false){10}else{20}", 20},
		{"if(1){10}", 10},
		{"if(1<2){10}", 10},
		{"if(1<2){10}else{20}", 10},
		{"if(1>2){10}else{20}", 20},
		{"if(1>2){10}", Null},
		{"if(false){10}", Null},
		{"!(if(false){10})", true},
		{"if((if(false){10})){10}else{20}", 20},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T){
	tests := []vmTestCase{
		{"let one=1; one", 1},
		{"let one=1; let two=2;one+two", 3},
		{"let one=1; let two=one+one;one+two", 3},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T){
	tests:=[]vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}