package parser

import (
	"fmt"
	"github.com/singlaanish56/Compiler-in-go/token"
	
)

func (p *Parser) Errors() []error{
	return p.errors
}

func (p *Parser) nextToken(){
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) checkPeek(tokenType token.TokenType) bool{
	if p.peekTokenIs(tokenType){
		p.nextToken()
		return true
	}
	
	p.peekError(tokenType)
	return false
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool{
	return p.peekToken.Type == tokenType
}

func (p *Parser) peekError(tokenType token.TokenType){
	err := fmt.Errorf("expected the next token to be %s, got %s", tokenType, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) currTokenIs(tokenType token.TokenType) bool{
	return p.currToken.Type == tokenType
}

func (p *Parser) addPrefix(tokenType token.TokenType, fn prefixParseFn){
	if _, exists := p.prefixParserMap[tokenType]; exists{
		err := fmt.Errorf("prefix function already exists for token type %s", tokenType)
		p.errors = append(p.errors, err)
		return
	}

	p.prefixParserMap[tokenType] = fn
}

func (p *Parser) addInfix(tokenType token.TokenType, fn infixParseFn){
	if _, exists := p.infixParserMap[tokenType]; exists{
		err := fmt.Errorf("infix function already exists for token type %s", tokenType)
		p.errors = append(p.errors, err)
		return
	}

	p.infixParserMap[tokenType] = fn
}

func (p *Parser) currentPrecedence() int{
	if p, ok := precendences[p.currToken.Type]; ok{
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int{
	if p, ok := precendences[p.peekToken.Type]; ok{
		return p
	}

	return LOWEST
}

var precendences = map[token.TokenType]int{
	token.PLUS: SUM,
}

const (
	_int = iota
	LOWEST
	UNDEFINED2
	UNDEFINED3
	SUM
	UNDEFINED5
	PREFIX
)