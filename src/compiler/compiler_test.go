package compiler

import (
	"fmt"
	"testing"

	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/lexer"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/parser"
)



type testCompilerStructs struct{
	input string
	expectedConstants[]interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T){
	tests :=[]testCompilerStructs{
		{"1+2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpAdd), code.Make(code.OpPop)}},
		{"1;2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpPop), code.Make(code.OpConstant, 1), code.Make(code.OpPop)}},
		{"1-2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpSub), code.Make(code.OpPop)}},
		{"1*2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpMul), code.Make(code.OpPop)}},
		{"2/1", []interface{}{2, 1}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpDiv), code.Make(code.OpPop)}},
		{"-1", []interface{}{1},[]code.Instructions{code.Make(code.OpConstant, 0),code.Make(code.OpMinus), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestBooleanArithmetic(t *testing.T){
	tests := []testCompilerStructs{
		{"true", []interface{}{}, []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpPop)}},
		{"false", []interface{}{}, []code.Instructions{code.Make(code.OpFalse), code.Make(code.OpPop)}},
		{"1>2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpGreaterThan), code.Make(code.OpPop)}},
		{"1<2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpLessThan), code.Make(code.OpPop)}},
		{"1==2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpEqual), code.Make(code.OpPop)}},
		{"1!=2", []interface{}{1, 2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpNotEqual), code.Make(code.OpPop)}},
		{"true==false", []interface{}{}, []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpFalse), code.Make(code.OpEqual), code.Make(code.OpPop)}},
		{"true!=false", []interface{}{}, []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpFalse), code.Make(code.OpNotEqual), code.Make(code.OpPop)}},
		{"!true", []interface{}{}, []code.Instructions{code.Make(code.OpTrue),code.Make(code.OpBang), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T){
	tests := []testCompilerStructs{
		{`if (true){10}; 3333;`, []interface{}{10, 3333}, []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpJumpNotTruthy, 10), code.Make(code.OpConstant, 0),code.Make(code.OpJump, 11), code.Make(code.OpNull) ,code.Make(code.OpPop), code.Make(code.OpConstant, 1), code.Make(code.OpPop)}},
		{`if (true){10}else{20}; 3333;`, []interface{}{10,20, 3333}, []code.Instructions{code.Make(code.OpTrue), code.Make(code.OpJumpNotTruthy, 10), code.Make(code.OpConstant, 0), code.Make(code.OpJump, 13), code.Make(code.OpConstant, 1), code.Make(code.OpPop), code.Make(code.OpConstant, 2), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestGlobalVariables(t *testing.T){
	tests := []testCompilerStructs{
		{`let one=1;let two=2;`,[]interface{}{1,2}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpSetGlobal, 0), code.Make(code.OpConstant, 1), code.Make(code.OpSetGlobal, 1)}},
		{`let one=1;one;`,[]interface{}{1}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpSetGlobal, 0), code.Make(code.OpGetGlobal, 0), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T){
	tests:= []testCompilerStructs{
		{`"monkey"`, []interface{}{"monkey"}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpPop)}},
		{`"mon" + "key"`, []interface{}{"mon","key"}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant,1),code.Make(code.OpAdd), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestArrayExpressions(t *testing.T){
	tests := []testCompilerStructs{
		{"[]",[]interface{}{}, []code.Instructions{code.Make(code.OpArray, 0), code.Make(code.OpPop)}},
		{"[1,2,3]",[]interface{}{1,2,3}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpConstant, 2), code.Make(code.OpArray, 3), code.Make(code.OpPop)}},
		{"[1+2,3-4, 5*6]",[]interface{}{1,2,3,4,5,6}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpAdd), code.Make(code.OpConstant,2), code.Make(code.OpConstant,3),code.Make(code.OpSub), code.Make(code.OpConstant,4), code.Make(code.OpConstant, 5), code.Make(code.OpMul), code.Make(code.OpArray,3), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestHashExpressions(t *testing.T){
	tests := []testCompilerStructs{
		{"{}",[]interface{}{}, []code.Instructions{code.Make(code.OpHash, 0), code.Make(code.OpPop)}},
		{"{1:2,3:4,5:6}", []interface{}{1,2,3,4,5,6}, []code.Instructions{code.Make(code.OpConstant, 0),code.Make(code.OpConstant, 1),code.Make(code.OpConstant, 2),code.Make(code.OpConstant, 3),code.Make(code.OpConstant, 4),code.Make(code.OpConstant, 5),code.Make(code.OpHash, 6), code.Make(code.OpPop)}},
		{"{1:2+3,4:5*6}", []interface{}{1,2,3,4,5,6},[]code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpConstant, 2), code.Make(code.OpAdd),code.Make(code.OpConstant,3), code.Make(code.OpConstant,4), code.Make(code.OpConstant, 5), code.Make(code.OpMul), code.Make(code.OpHash,4), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T){
	tests := []testCompilerStructs{
		{"[1,2,3][1+1]", []interface{}{1,2,3,1,1}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpConstant,2), code.Make(code.OpArray,3), code.Make(code.OpConstant, 3), code.Make(code.OpConstant, 4), code.Make(code.OpAdd), code.Make(code.OpIndex), code.Make(code.OpPop)}},
		{"{1:2}[2-1]", []interface{}{1,2,2,1}, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpConstant, 1), code.Make(code.OpHash,2) , code.Make(code.OpConstant, 2), code.Make(code.OpConstant, 3), code.Make(code.OpSub), code.Make(code.OpIndex), code.Make(code.OpPop)}},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T){
	tests := []testCompilerStructs{
		{
			`fn(){return 5+10;}`,
			[]interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			`fn(){5+10;}`,
			[]interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			`fn(){1; 2}`,
			[]interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			[]code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			`fn(){}`,
			[]interface{}{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			[]code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T){
	compiler := New()

	if compiler.scopeIndex != 0{
		t.Errorf("scopeIndex wrong, expected=0, got=%d", compiler.scopeIndex)
	}
	globalSymbolTable  := compiler.symbolTable
	compiler.emit(code.OpMul)

	compiler.enterScope()

	if compiler.scopeIndex != 1{
		t.Errorf("scopeIndex wrong, expected=1, got=%d", compiler.scopeIndex)
	}

	compiler.emit(code.OpAdd)

	if len(compiler.compilerScopes[compiler.scopeIndex].instructions) != 1{
		t.Errorf("instructions length wrong, expected=1, got=%d", len(compiler.compilerScopes[compiler.scopeIndex].instructions))
	}

	last := compiler.compilerScopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpAdd{
		t.Errorf("last instruction wrong, expected=%d, got=%d", code.OpAdd, last.Opcode)
	}

	if compiler.symbolTable.Outer != globalSymbolTable{
		t.Errorf("compiler doesnt include the enclose symbopl table")
	}

	compiler.leaveScope()

	if compiler.scopeIndex != 0{
		t.Errorf("scopeIndex wrong, expected=0, got=%d", compiler.scopeIndex)
	}

	if compiler.symbolTable != globalSymbolTable{
		t.Errorf("compiler doesnt include the global symbol table")
	}
	if compiler.symbolTable.Outer != nil{
		t.Errorf("compiler symbol table should not have outer symbol table")
	}

	compiler.emit(code.OpSub)

	if len(compiler.compilerScopes[compiler.scopeIndex].instructions) != 2{
		t.Errorf("instructions length wrong, expected=2, got=%d", len(compiler.compilerScopes[compiler.scopeIndex].instructions))
	}

	last = compiler.compilerScopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpSub{
		t.Errorf("last instruction wrong, expected=%d, got=%d", code.OpSub, last.Opcode)
	}

	prev := compiler.compilerScopes[compiler.scopeIndex].previousInstruction
	if prev.Opcode != code.OpMul{
		t.Errorf("previous instruction wrong, expected=%d, got=%d", code.OpMul, prev.Opcode)
	}
}

func TestFunctionCalls(t *testing.T){
tests := []testCompilerStructs{
{
	`fn(){24}();`,
	[]interface{}{
		24,
		[]code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpReturnValue),
		},
	},
	[]code.Instructions{
		code.Make(code.OpConstant, 1),
		code.Make(code.OpCall,0),
		code.Make(code.OpPop),
	},
},
{
	`let callArg = fn(){24};
	callArg();`,
	[]interface{}{
		24,
		[]code.Instructions{
			code.Make(code.OpConstant, 0),
			code.Make(code.OpReturnValue),
		},
	},
	[]code.Instructions{
		code.Make(code.OpConstant, 1),
		code.Make(code.OpSetGlobal, 0),
		code.Make(code.OpGetGlobal, 0),
		code.Make(code.OpCall,0),
		code.Make(code.OpPop),
	},
},
{
	`let oneArg = fn(a){a}; oneArg(24);`,
	[]interface{}{
		[]code.Instructions{
			code.Make(code.OpGetLocal, 0),
			code.Make(code.OpReturnValue),
		},
		24,
	},
	[]code.Instructions{
		code.Make(code.OpConstant, 0),
		code.Make(code.OpSetGlobal, 0),
		code.Make(code.OpGetGlobal, 0),
		code.Make(code.OpConstant, 1),
		code.Make(code.OpCall, 1),
		code.Make(code.OpPop),
	},
},
{
	`let manyArg = fn(a, b, c){a; b; c;}; manyArg(24, 48 , 36);`,
	[]interface{}{
		[]code.Instructions{
			code.Make(code.OpGetLocal, 0),
			code.Make(code.OpPop),
			code.Make(code.OpGetLocal, 1),
			code.Make(code.OpPop),
			code.Make(code.OpGetLocal, 2),
			code.Make(code.OpReturnValue),
		},
		24,
		48,
		36,
	},
	[]code.Instructions{
		code.Make(code.OpConstant, 0),
		code.Make(code.OpSetGlobal, 0),
		code.Make(code.OpGetGlobal, 0),
		code.Make(code.OpConstant, 1),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 3),
		code.Make(code.OpCall, 3),
		code.Make(code.OpPop),
	},
},
}	

runCompilerTests(t, tests)
}

func TestLetStatementScopes(t *testing.T){
	tests := []testCompilerStructs{
		{`let num=55; fn(){num}`, 
		[]interface{}{55, []code.Instructions{code.Make(code.OpGetGlobal, 0), code.Make(code.OpReturnValue)}},
		[]code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpSetGlobal, 0), code.Make(code.OpConstant, 1), code.Make(code.OpPop)},
		},
		{`fn(){let num=55;  num}`, 
		[]interface{}{55, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpSetLocal, 0), code.Make(code.OpGetLocal, 0), code.Make(code.OpReturnValue)}},
		[]code.Instructions{code.Make(code.OpConstant, 1), code.Make(code.OpPop)},
		},
		{`fn(){let a=55;let b= 77;  a+b}`, 
		[]interface{}{55,77, []code.Instructions{code.Make(code.OpConstant, 0), code.Make(code.OpSetLocal, 0),code.Make(code.OpConstant, 1), code.Make(code.OpSetLocal, 1), code.Make(code.OpGetLocal, 0), code.Make(code.OpGetLocal, 1), code.Make(code.OpAdd), code.Make(code.OpReturnValue)}},
		[]code.Instructions{code.Make(code.OpConstant, 2), code.Make(code.OpPop)},
		},
	}

	runCompilerTests(t, tests)
}



func runCompilerTests(t *testing.T, tests []testCompilerStructs){
	t.Helper()

	for _, tt := range tests{
		program := parse(tt.input)

		compiler := New()

		err := compiler.Compile(program)

		if err != nil{
			t.Fatalf("compiler error %s", err)
		}

		bytecode := compiler.Bytecode()
		err = testInstructions(bytecode.Instructions, tt.expectedInstructions)
		if err != nil{
			t.Fatalf("instructions dont match %s", err)
		}

		err= testConstants(bytecode.Constants, tt.expectedConstants)
		if err != nil{
			t.Fatalf("constants dont match %s", err)
		}


	}
}

func parse(input string) *ast.AstRootNode{
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParserProgram()
}

func testInstructions(actual code.Instructions, expected []code.Instructions) error{

	concatted := concatInstructions(expected)

	if len(actual) != len(concatted){
		return fmt.Errorf("wrong instruction length, expected=%q, got=%q", concatted, actual)
	}

	for i, ins := range concatted{
		if actual[i] != ins{
			return fmt.Errorf("wrong instruction at %d, expected=%q, got=%q", i, concatted, actual)
		}
	}

	return nil
}

func testConstants(actual []object.Object, expected []interface{}) error{
	if len(actual) != len(expected){
		return fmt.Errorf("wrong number of constants, expected=%d, got=%d", len(expected), len(actual))
	}

	for i, constant := range expected{
		switch constant := constant.(type){
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil{
				return fmt.Errorf("constant at index %d, expected=%d, got=%s", i, constant, err)	
			}
		case string:
			err := testStringObject(constant, actual[i])
			if err != nil{
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		case []code.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok{
				return fmt.Errorf("object is not a function, got=%T", actual[i])
			}

			err := testInstructions(fn.Instructions, constant)
			if err != nil{
				return fmt.Errorf("constant %d - testInstructions failed: %s", i, err)
			}
		}
	}

	return nil
}

func concatInstructions(expected []code.Instructions) code.Instructions{
	var concatted code.Instructions

	for _, ins := range expected{
		concatted = append(concatted, ins...)
	}

	return concatted
}

func testIntegerObject(expected int64, actual object.Object) error{
	result, ok := actual.(*object.Integer)

	if !ok{
		return fmt.Errorf("object is not Integer, got=%T", actual)
	}

	if result.Value != expected{
		return fmt.Errorf("object has wrong value, expected=%d, got=%d", expected, result.Value)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error{
	result, ok := actual.(*object.String)
	if !ok{
		return fmt.Errorf("the error is not a string, got=%T", actual)
	}

	if result.Value != expected{
		return fmt.Errorf("object has the wrong value, expected=%s, got=%s", expected, result.Value)
	}
	return nil
}