package ast

import (
	"bytes"
	"github.com/singlaanish56/Compiler-in-go/token"
)

type ASTNode interface{
	TokenLiteral() string
	String() string
}

type Statement interface{
	ASTNode
	statementNode()
}

type Expression interface{
	ASTNode
	expressionNode()
}

type AstRootNode struct{
	Statements []Statement
}

func (root *AstRootNode) TokenLiteral() string{
	if len(root.Statements)>0{
		return root.Statements[0].TokenLiteral()
	}
	return ""
}

func (root *AstRootNode) String() string{
	var out bytes.Buffer
	for _, stmt := range root.Statements{
		out.WriteString(stmt.String())
	}

	return out.String()
}

type LetStatement struct{
	Token token.Token
	Variable *Variable
	Value Expression
}
func (ls *LetStatement) statementNode(){}
func (ls *LetStatement) TokenLiteral() string{
	return ls.Token.Identifier
}
func (ls *LetStatement) String() string{
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Variable.String())
	out.WriteString(" = ")
	if ls.Value != nil{
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct{
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode(){}
func (rs *ReturnStatement) TokenLiteral() string{return rs.Token.Identifier}
func (rs *ReturnStatement) String() string{
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil{
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct{
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode(){}
func (es *ExpressionStatement) TokenLiteral() string{return es.Token.Identifier}
func (es *ExpressionStatement) String() string{
	
	if es.Expression != nil{
	return es.Expression.String()	
	}
	return ""
}

type Variable struct{
	Token token.Token
	Value string
}

func (v *Variable) expressionNode(){}
func (v *Variable) TokenLiteral() string{return v.Token.Identifier}
func (v *Variable) String() string{ return v.Value}

type IntegerLiteral struct{
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode(){}
func (il *IntegerLiteral) TokenLiteral() string{return il.Token.Identifier}
func (il *IntegerLiteral) String() string{return il.Token.Identifier}

type BooleanLiteral struct{
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode(){}
func (bl *BooleanLiteral) TokenLiteral() string{return bl.Token.Identifier}
func (bl *BooleanLiteral) String() string{return bl.Token.Identifier}

type StringLiteral struct{
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode(){}
func (sl *StringLiteral) TokenLiteral() string{return sl.Token.Identifier}
func (sl *StringLiteral) String() string{return sl.Token.Identifier}

type PrefixExpression struct{
	Token token.Token
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode(){}
func (pe *PrefixExpression) TokenLiteral() string{return pe.Token.Identifier}
func (pe *PrefixExpression) String() string{
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	if pe.Right != nil{
		out.WriteString(pe.Right.String())
	}
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct{
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (ie *InfixExpression) expressionNode(){}
func (ie *InfixExpression) TokenLiteral() string{return ie.Token.Identifier}
func (ie *InfixExpression) String() string{
	var out bytes.Buffer
	out.WriteString("(")
	if ie.Left != nil{
		out.WriteString(ie.Left.String())
	}
	out.WriteString(" " + ie.Operator + " ")
	if ie.Right != nil{
		out.WriteString(ie.Right.String())
	}
	out.WriteString(")")
	return out.String()
}