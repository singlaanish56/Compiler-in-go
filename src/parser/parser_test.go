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

func TestArrayLiteral(t *testing.T){
	input := "[1, 2*2, 3+3]"

	l := lexer.New(input)
	p := New(l)
	prog := p.ParserProgram()

	array , ok := prog.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.ArrayLiteral)
	if !ok{
		t.Fatalf("the type of the array is not array literal , got=%T", prog.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if len(array.Elements) != 3{
		t.Fatalf("the size of the array not as expected= 3, got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfix(t, array.Elements[1], 2, "*", 2)
	testInfix(t, array.Elements[2], 3, "+", 3)
}

func TestParserIndexExpression(t *testing.T){
	input := "arr[1+1]"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParserProgram()
	
	if len(prog.Statements) != 1{
		t.Errorf("the number of statements not as expected")
		return 
	}

	st, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Errorf("the expression type not as expected, got=%T", prog.Statements[0])
		return 
	}

	exp, ok := st.Expression.(*ast.IndexExpression)
	if !ok{
		t.Errorf("the expression type not as expected, got=%T", prog.Statements[0])
		return
	}

	if !testIdentifier(t, exp.Left, "arr"){
		return 
	}

	if !testInfix(t, exp.Index, 1, "+", 1){
		return 
	}
}

func TestHashLiteral(t *testing.T){
 input := `{"one":1,"second":2,"third":3}`

 l := lexer.New(input)
 p := New(l)
 prog := p.ParserProgram()
 
 if len(prog.Statements) != 1{
	t.Errorf("the number of statements is not as expected")
	return
 }

 st, ok := prog.Statements[0].(*ast.ExpressionStatement)
 if !ok{
	t.Errorf("the expression type is not expected got=%T", prog.Statements[0])
	return
 }

 hash, ok := st.Expression.(*ast.HashLiteral)
 if !ok{
	t.Errorf("the number of key value pairs is not as expected got=%T",st.Expression)
	return 
 }
 if len(hash.Pairs) != 3{
	t.Errorf("the number of kv pairss not as expected , expected=3, got=%d", len(hash.Pairs))
	return
 }

 expected := map[string]int{
	"one":1,
	"second":2,
	"third":3,
 }
 for k, v := range hash.Pairs{
	literal, ok := k.(*ast.StringLiteral)
	if !ok{
		t.Errorf("the key is not a string literal , got=%T", literal)
	}

	expectedValue := expected[literal.String()]

	testIntegerLiteral(t, v, int64(expectedValue))
 }
}

func  TestIfExpression(t *testing.T){
	tests := []struct{
		input string
	}{
		{"if(x<y){x}"},
		{"if(x<y){x}else{y}"},
	}

	for _, tt := range tests{
	
	l := lexer.New(tt.input)
	p := New(l)
	prog := p.ParserProgram()

	if len(prog.Statements) !=1 {
		t.Errorf("the len of the statements  is wrong, got=%d", len(prog.Statements))
		return 
	}

	st, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Errorf("the program statement is of wrong type , got=%T", prog.Statements[0])
		return
	}

	exp, ok := st.Expression.(*ast.IfExpression)
	if !ok{
		t.Errorf("expected the if expression , got=%T",st.Expression)
		return
	}

	if !testInfix(t, exp.Condition, "x", "<", "y"){
		return 
	}

	if len(exp.Consequence.Statements) != 1{
		t.Errorf("the consquence is not enoughs statements")
		return 
	}

	con, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok{
		t.Errorf("the consquence the if doesnt have the expected type got=%T", con)
		return 
	}
	if !testIdentifier(t, con.Expression, "x"){
		return 
	}
	if exp.Alternative != nil{
		if len(exp.Alternative.Statements) != 1{
			t.Errorf("the consquence has not enough statements")
			return 
		}
		alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Errorf("the consquence for the if doesnt have the expected type=%T",exp.Alternative.Statements[0])
			return 
		}
		if !testIdentifier(t, alt.Expression, "y"){
			return 
		}
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

func testInfix(t *testing.T, exp ast.Expression, left interface{}, operator string , right interface{}) bool{
inexp, ok := exp.(*ast.InfixExpression)
if !ok{
	t.Errorf("exp not infix, got=%T", inexp)
	return false
}

if !testLiteralExpression(t, inexp.Left, left){
	return false
}

if inexp.Operator != operator{
	t.Errorf("operator not matchinf for the infix expression, got=%q, expected=%q", inexp.Operator, operator)
	return false
}

if !testLiteralExpression(t, inexp.Right, right){
	return false
}

return true
}