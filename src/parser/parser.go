package parser

import (
	"fmt"

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
	p.infixParserMap = make(map[token.TokenType]infixParseFn)
	p.addPrefix(token.VARIABLE, p.parseVariable)
	p.addPrefix(token.NUMBER, p.parseNumber)
	p.addPrefix(token.STRING, p.parseStringExpression)

	p.addPrefix(token.MINUS, p.parsePrefixExpression)
	p.addPrefix(token.PLUS, p.parsePrefixExpression)
	p.addPrefix(token.EXCLAMATION, p.parsePrefixExpression)

	p.addPrefix(token.TRUE, p.parseBooleanExpression)
	p.addPrefix(token.FALSE, p.parseBooleanExpression)


	p.addInfix(token.PLUS, p.parseInfixExpression)
	p.addInfix(token.MINUS, p.parseInfixExpression)

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
	return rootNode
}

func (p* Parser) parseStatement() ast.Statement{
	switch p.currToken.Type{
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() ast.Statement{
	st := &ast.ExpressionStatement{Token: p.currToken}
	st.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return st
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



