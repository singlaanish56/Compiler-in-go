package ast

import (
	"testing"

	"github.com/singlaanish56/Compiler-in-go/token"
)

func TestString(t *testing.T){
	rootNode := &AstRootNode{
		Statements: []Statement{
			&LetStatement{
				Token : token.Token{Type: token.LET, Identifier: "let"},
				Variable: &Variable{Token: token.Token{Type: token.VARIABLE, Identifier: "x"},Value: "x"},
				Value: &Variable{Token: token.Token{Type: token.VARIABLE, Identifier: "another one"},Value: "another one"},
			},
		},
	}

	if rootNode.String() != "let x = another one;"{
		t.Errorf("Expected 'let x = another one;', got '%s'", rootNode.String())
	}
}