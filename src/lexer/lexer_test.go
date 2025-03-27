package lexer

import (
	"fmt"
	"testing"

	"github.com/singlaanish56/Compiler-in-go/token"
)

func TestNextToken(t * testing.T){
	input := `let five=5;let ten=10;let add = fn(x,y){x+y;}let str = "this is a string";!-/*5;5<10>5; if(5<10){return true;}else{return false;}10==10;10!=9;[1,2];:`

	tests:= []struct{
		expectedType token.TokenType
		expectedIdentifier string
	}{
		{token.LET,"let"},
		{token.VARIABLE,"five"},
		{token.EQUALTO,"="},
		{token.NUMBER,"5"},
		{token.SEMICOLON,";"},
		{token.LET,"let"},
		{token.VARIABLE,"ten"},
		{token.EQUALTO,"="},
		{token.NUMBER,"10"},
		{token.SEMICOLON,";"},
		{token.LET,"let"},
		{token.VARIABLE,"add"},
		{token.EQUALTO,"="},
		{token.FUNCTION,"fn"},
		{token.OPENROUND,"("},
		{token.VARIABLE,"x"},
		{token.COMMA,","},
		{token.VARIABLE,"y"},
		{token.CLOSEROUND,")"},
		{token.OPENBRACE,"{"},
		{token.VARIABLE,"x"},
		{token.PLUS,"+"},
		{token.VARIABLE,"y"},
		{token.SEMICOLON,";"},
		{token.CLOSEBRACE,"}"},
		{token.LET,"let"},
		{token.VARIABLE,"str"},
		{token.EQUALTO,"="},
		{token.STRING,"this is a string"},
		{token.SEMICOLON,";"},
		{token.EXCLAMATION,"!"},
		{token.MINUS,"-"},
		{token.DIVIDE,"/"},
		{token.MULTIPLY,"*"},
		{token.NUMBER,"5"},
		{token.SEMICOLON,";"},
		{token.NUMBER,"5"},
		{token.OPENANGLE,"<"},
		{token.NUMBER,"10"},
		{token.CLOSEANGLE,">"},
		{token.NUMBER,"5"},
		{token.SEMICOLON,";"},
		{token.IF,"if"},
		{token.OPENROUND,"("},
		{token.NUMBER,"5"},
		{token.OPENANGLE,"<"},
		{token.NUMBER,"10"},
		{token.CLOSEROUND,")"},
		{token.OPENBRACE,"{"},
		{token.RETURN,"return"},
		{token.TRUE,"true"},
		{token.SEMICOLON,";"},
		{token.CLOSEBRACE,"}"},
		{token.ELSE,"else"},
		{token.OPENBRACE,"{"},
		{token.RETURN,"return"},
		{token.FALSE,"false"},
		{token.SEMICOLON,";"},
		{token.CLOSEBRACE,"}"},
		{token.NUMBER,"10"},
		{token.DOUBLEEQUALTO,"=="},
		{token.NUMBER,"10"},
		{token.SEMICOLON,";"},
		{token.NUMBER,"10"},
		{token.EXCLAMATIONEQUALTO,"!="},
		{token.NUMBER,"9"},
		{token.SEMICOLON,";"},
		{token.OPENBRACKET,"["},
		{token.NUMBER,"1"},
		{token.COMMA,","},
		{token.NUMBER,"2"},
		{token.CLOSEBRACKET,"]"},
		{token.SEMICOLON,";"},
		{token.COLON,":"},
		{token.EOF,""},	
	}

	l := New(input)

	for i, tt := range tests{
		tok := l.NextToken()

		if tok.Type != tt.expectedType{
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Identifier != tt.expectedIdentifier{
			t.Fatalf("tests[%d] - tokenIdentifier wrong. expected=%q, got=%q", i, tt.expectedIdentifier, tok.Identifier)
		}

		fmt.Printf("tokenLiteral : %q\n", tt.expectedIdentifier)
	}
}