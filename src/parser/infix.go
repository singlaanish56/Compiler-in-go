package parser

import (
	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/token"
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

func (p *Parser) parserArrayIndexExpression(left ast.Expression) ast.Expression{
	indexExp := &ast.IndexExpression{Token: p.currToken, Left: left}

	p.nextToken()

	indexExp.Index = p.parseExpression(LOWEST)

	if !p.checkPeek(token.CLOSEBRACKET){
		return nil
	}

	return indexExp
}