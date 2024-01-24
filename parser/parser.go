package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
	"fmt"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{}}

	//Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {

	//Constucts the root node of the AST
	program := &ast.Program{}
	//Initialises the Statements part
	program.Statements = []ast.Statement{}

	//Iterates over every token until it reaches an EOF
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		//If stmt is anything other than 'nil', it appends it to 'program.statements'
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// This function is in charge of parsing all the different statements
func (p *Parser) parseStatement() ast.Statement {
	//Get the current tokens type
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// This function parses all the LET statements
func (p *Parser) parseLetStatement() *ast.LetStatement {
	//Constructs a LetStatement node with the token it's sitting on (a token.LET token)
	stmt := &ast.LetStatement{Token: p.curToken}

	//If the next token ISN'T a variable name, it returns with nil (it's no longer)
	// (a let statement)
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	//Since the next token is a variable name:
	//It sets the Let Statements name to the same as the variable name
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	//We then need an equals sign for the let statement
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
