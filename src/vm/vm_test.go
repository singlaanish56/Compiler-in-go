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
	case []int:
		array, ok := obj.(*object.Array)
		if !ok{
			t.Errorf("object is not an Array, got=%T(%+v)", obj, obj)
		}

		if len(array.Elements) != len(expected){
			t.Errorf("wrong number of elements in array, expected=%d, got=%d", len(expected), len(array.Elements))
		}

		for i, expectedElem := range expected{
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil{
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := obj.(*object.Hash)
		if !ok{
			t.Errorf("object is not a Hash, got=%T(%+v)", obj, obj)
			return 
		}
		if len(hash.Pairs) != len(expected){
			t.Errorf("wrong number of pairs in hash, expected=%d, got=%d", len(expected), len(hash.Pairs))
			return
		}
		for expectedKey, expectedValue := range expected{
			pair, ok := hash.Pairs[expectedKey]
			if !ok{
				t.Errorf("no pair for given key")
			}
			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil{
				t.Errorf("testIntegerObject failed: %s", err)
			}
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

func TestArrayExpressions(t *testing.T){
	tests:=[]vmTestCase{
		{"[]",[]int{}},
		{"[1,2,3]",[]int{1,2,3}},
		{"[1+2, 4-5, 6*8]",[]int{3,-1,48}},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T){
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{"{1:2,3:4}", map[object.HashKey]int64{
			(&object.Integer{Value:1}).HashKey(): 2,
			(&object.Integer{Value:3}).HashKey(): 4,
			},
		},
		{"{1+1:2*2, 3+3:4*4}", map[object.HashKey]int64{
			(&object.Integer{Value:2}).HashKey(): 4,
			(&object.Integer{Value:6}).HashKey(): 16,	
			},
		},
	}

	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T){
	tests := []vmTestCase{
		{"[1,2,3][1]",2},
		{"[1,2,3][0+2]",3},
		{"[[1,1,1]][0][0]",1},
		{"[][0]",Null},
		{"[1,2,3][99]",Null},
		{"[1][-1]",Null},
		{"{1:1, 2:2}[1]",1},
		{"{1:1, 2:2}[2]",2},
		{"{}[0]",Null},
	}

	runVmTests(t, tests)
}

func TestFunctions(t *testing.T){
	tests := []vmTestCase{
		{`let fivePlusTen = fn(){5+10;}; fivePlusTen();`, 15},
		{`let one = fn(){1;}; let two=fn(){2;}; one() + two()`, 3},
		{`let a = fn(){1}; let b=fn(){a()+1}; let c=fn(){b()+1}; c();`, 3},
		{`let a = fn(){return 99; 100;}; a();`, 99},
		{`let a = fn(){return 99; return 100;}; a();`, 99},
		{`let a = fn(){}; a();`, Null},
		{`let a = fn(){}; let b = fn(){a();}; a(); b();`, Null},
		{`let a = fn(){1;}; let b = fn(){a;}; b()();`, 1},
		{`let a = fn(){let one =1;one}; a();`, 1},
		{`let a = fn(){let one =1; let two=2; one+two}; a();`, 3},
		{`let a = fn(){let one =1; let two=2; one+two}; let b =fn(){let one=3; let two=4; one+two}; a() + b();`, 10},
		{`let global=50;let minus= fn(){let num=1; global-num;}; let plus=fn(){let num=1; global+num;}; minus() + plus();`, 100},
		{`let returnsOneReturner = fn(){let returnsOne= fn(){1;}; returnsOne();}; returnsOneReturner();`, 1},
	}

	runVmTests(t, tests)
}