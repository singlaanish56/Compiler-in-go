package parser

import (
	"fmt"
	"strconv"

	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/token"
	
)

func (p *Parser) parseVariable() ast.Expression{
	return &ast.Variable{Token: p.currToken, Value: p.currToken.Identifier}
}

func (p *Parser) parseNumber() ast.Expression{
	integerLiteral := &ast.IntegerLiteral{Token: p.currToken}

	val, err := strconv.ParseInt(p.currToken.Identifier, 0 , 64)
	if err != nil{
		p.errors = append(p.errors, fmt.Errorf("could not parse %s as integer", p.currToken.Identifier))
		return nil
	}

	integerLiteral.Value = val
	return integerLiteral
}

func (p *Parser) parsePrefixExpression() ast.Expression{
	prefixExpression := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Identifier}
	p.nextToken()
	prefixExpression.Right = p.parseExpression(PREFIX)

	return prefixExpression
}

func (p *Parser) parseStringExpression() ast.Expression{
	strLiteral := &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Identifier}
	return strLiteral
}

func (p *Parser) parseBooleanExpression() ast.Expression{
	return &ast.BooleanLiteral{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}