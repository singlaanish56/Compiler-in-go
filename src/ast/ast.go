package ast

import (
	"bytes"
	"strings"

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

type BlockStatement struct{
	Token token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode(){}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Identifier}
func (bs *BlockStatement) String() string{
	var out bytes.Buffer
	for _, s := range bs.Statements{
		out.WriteString(s.String())
	}

	return out.String()
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

type ArrayLiteral struct{
	Token token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode(){}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Identifier}
func (al *ArrayLiteral) String() string{
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements{
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements,","))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct{
	Token token.Token
	Left Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode(){}
func (ie *IndexExpression) TokenLiteral() string{ return ie.Token.Identifier}
func (ie *IndexExpression) String() string{
  var out bytes.Buffer

  out.WriteString("(")
  out.WriteString(ie.Left.String())
  out.WriteString("[")
  out.WriteString(ie.Index.String())
  out.WriteString("]")
  out.WriteString(")")

  return out.String()
}

type HashLiteral struct{
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode(){}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Identifier}
func (hl *HashLiteral) String() string{
	var out bytes.Buffer

	elements := []string{}

	for k, v := range hl.Pairs{
		elements = append(elements, k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements,","))
	out.WriteString("}")

	return out.String()
}

type IfExpression struct{
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode(){}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Identifier}
func (ie *IfExpression) String() string{
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil{
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

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