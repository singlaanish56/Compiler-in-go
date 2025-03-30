package parser

import (
	"fmt"
	"testing"

	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/lexer"

)

func TestLetStatement(t *testing.T){
	tests := []struct{
		input string
		expectIdentifier string
		expectedValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;","foobar", "y"},
	}

	for _, tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParserProgram()
		if len(p.Errors()) != 0{
			t.Errorf("Parser has errors: %v", p.Errors())
		}

		if len(program.Statements) != 1{
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		st := program.Statements[0]
		if !testLetStatement(t, st, tt.expectIdentifier){
			return 
		}

		val := st.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue){
			return 
		}
	}
}

func TestReturnStatement(t *testing.T){
	tests := []struct{
		input string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return foobar;", "foobar"},
		{"return false;", false},
	}

	for _, tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParserProgram()

		if len(program.Statements) != 1{
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		st := program.Statements[0]
		val := st.(*ast.ReturnStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue){
			return 
		}
	}
}

//helper functions

func testLetStatement(t *testing.T, s ast.Statement, name string) bool{
	if s.TokenLiteral() != "let"{
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letst, ok := s.(*ast.LetStatement)
	if !ok{
		t.Errorf("this is statement is not a let statement")
		return false
	}
	if letst.Variable.Value != name{
		t.Errorf("letst.Variable.Value not %s. got=%s", name, letst.Variable.Value)
		return false
	}

	if letst.Variable.TokenLiteral() != name{
		t.Errorf("letst.Variable.TokenLiteral not %s. got=%s", name, letst.Variable.TokenLiteral())
		return false
	}	

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool{
	switch v := expected.(type){
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
		return false
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, expected int64) bool{
	iliteral, ok := il.(*ast.IntegerLiteral)
	if !ok{
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if iliteral.Value != expected{
		t.Errorf("iliteral.Value not %d. got=%d", expected, iliteral.Value)
		return false
	}

	if iliteral.TokenLiteral() != fmt.Sprintf("%d", expected){
		t.Errorf("iliteral.TokenLiteral not %d. got=%s", expected, iliteral.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, sl ast.Expression, expected string) bool{
	sliteral, ok := sl.(*ast.Variable)
	if !ok{
		t.Errorf("sl not *ast.StringLiteral. got=%T", sl)
		return false
	}

	if sliteral.Value != expected{
		t.Errorf("sliteral.Value not %s. got=%s", expected, sliteral.Value)
		return false
	}

	if sliteral.TokenLiteral() != expected{
		t.Errorf("sliteral.TokenLiteral not %s. got=%s", expected, sliteral.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, b ast.Expression, expected interface{}) bool{
	bo , ok := b.(*ast.BooleanLiteral)
	if !ok{
		t.Errorf("the expression is not expected, got=%T", b)
		return false
	}

	if bo.Value != expected{
		t.Errorf("the value is not expected=%t, got=%t", expected, bo.Value)
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", expected){
		t.Errorf("the value is not expected=%q, got=%q", bo.TokenLiteral(),fmt.Sprintf("%t", expected))
		return false
	} 

	return true
}