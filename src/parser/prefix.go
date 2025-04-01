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

func (p *Parser) parseArrayExpression() ast.Expression{
	arr :=  &ast.ArrayLiteral{Token: p.currToken}

	arr.Elements = p.parseExpressionList(token.CLOSEBRACKET)

	return arr
}

func (p *Parser) parseGroupedExpression() ast.Expression{
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.checkPeek(token.CLOSEROUND){
		return nil
	}

	return exp
}

func (p *Parser) parseHashMapExpression() ast.Expression{
	hashExp := &ast.HashLiteral{Token:p.currToken}
	hashExp.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.CLOSEBRACE){
		p.nextToken()

		key:=p.parseExpression(LOWEST)
			
		if !p.checkPeek(token.COLON){
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hashExp.Pairs[key] = value

		if !p.peekTokenIs(token.CLOSEBRACE) && !p.checkPeek(token.COMMA){
			return nil
		}
	}

	if !p.checkPeek(token.CLOSEBRACE){
		return nil
	}

	return hashExp
}

func (p *Parser) parseIfExpression() ast.Expression{
	exp := &ast.IfExpression{Token: p.currToken}

	if !p.checkPeek(token.OPENROUND){
		return nil
	}

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)


	if !p.checkPeek(token.CLOSEROUND){
		return nil
	}

	if !p.checkPeek(token.OPENBRACE){
		return nil
	}
	
	exp.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE){
		p.nextToken()

		if !p.checkPeek(token.OPENBRACE){
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression{
	exp := &ast.FunctionExpression{Token: p.currToken}

	if !p.checkPeek(token.OPENROUND){
		return nil
	}

	exp.Parameters = p.parseFunctionArguments()

	if !p.checkPeek(token.OPENBRACE){
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFunctionArguments() []*ast.Variable{
	params := []*ast.Variable{}

	if p.peekTokenIs(token.CLOSEROUND){
		p.nextToken()
		return params
	}

	p.nextToken()

	arg := &ast.Variable{Token: p.currToken, Value: p.currToken.Identifier}
	params = append(params, arg)

	for p.peekTokenIs(token.COMMA){
		p.nextToken()
		p.nextToken()

		arg = &ast.Variable{Token: p.currToken, Value: p.currToken.Identifier}
		params = append(params, arg)
	}

	if !p.checkPeek(token.CLOSEROUND){
		return nil
	}

	return params
}