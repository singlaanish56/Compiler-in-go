package parser

import (
	"fmt"
	"strconv"

	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/lexer"
	"github.com/singlaanish56/Compiler-in-go/token"
	
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct{
	currToken token.Token
	peekToken token.Token
	lexer *lexer.Lexer
	errors []error

	prefixParserMap map[token.TokenType]prefixParseFn
	infixParserMap map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser{
	p := &Parser{lexer: lexer, errors: []error{}}
	p.nextToken()//sets the current token
	p.nextToken()//sets the next token

	p.prefixParserMap = make(map[token.TokenType]prefixParseFn)
	p.addPrefix(token.VARIABLE, p.parseVariable)
	p.addPrefix(token.NUMBER, p.parseNumber)
	p.addPrefix(token.PLUS, p.parsePrefixExpression)

	p.addInfix(token.PLUS, p.parseInfixExpression)
	return p
}

func (p *Parser) ParserProgram() *ast.AstRootNode{
	rootNode := &ast.AstRootNode{Statements: []ast.Statement{}}
	for p.currToken.Type != token.EOF{
		stmt := p.parseStatement()
		if stmt != nil{
			rootNode.Statements = append(rootNode.Statements, stmt)
		}

		p.nextToken()
	}

}

func (p *Parser) Errors() []error{
	return p.errors
}

func (p* Parser) parseStatement() ast.Statement{
	switch p.currToken.Type{
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}
 
func (p *Parser) parseLetStatement() ast.Statement{
	letstmt := & ast.LetStatement{Token: p.currToken}

	if !p.checkPeek(token.VARIABLE){
		return nil
	}

	letstmt.Variable = &ast.Variable{Token: p.currToken, Value:p.currToken.Identifier}

	if !p.checkPeek(token.EQUALTO){
		return nil
	}

	p.nextToken()

	letstmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	return letstmt
}

func (p *Parser) parseReturnStatement() ast.Statement{
	returnstmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	returnstmt.Value = p.parseExpression(LOWEST)

	for !p.currTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	return returnstmt
}

func (p *Parser) parseExpression(precendence int) ast.Expression{
	 prefixFn := p.prefixParserMap[p.currToken.Type]
	 if prefixFn == nil{
		p.errors = append(p.errors, fmt.Errorf("no prefix parse function for %s found", p.currToken.Type))
		return nil
	 }

	 leftExpression := prefixFn()

	 for !p.peekTokenIs(token.SEMICOLON) && precendence < p.peekPrecedence(){
		infixFn := p.infixParserMap[p.peekToken.Type]
		if infixFn == nil{
			return leftExpression
		}
		p.nextToken()
		leftExpression = infixFn(leftExpression)
	 }

	 return leftExpression

}

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