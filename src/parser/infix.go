package parser

import (

	"github.com/singlaanish56/Compiler-in-go/ast"
	
)


func (p *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression{
	exp := &ast.InfixExpression{
		Token: p.currToken,
		Operator: p.currToken.Identifier,
		Left: leftExpression,
	}

	precendence := p.currentPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precendence)
	return exp
}