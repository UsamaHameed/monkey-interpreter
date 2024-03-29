package parser

import (
	"fmt"
	"strconv"

	"github.com/UsamaHameed/monkey-interpreter/ast"
	"github.com/UsamaHameed/monkey-interpreter/lexer"
	"github.com/UsamaHameed/monkey-interpreter/token"
)

type (
    prefixParseFn   func() ast.Expression
    infixParseFn    func(ast.Expression) ast.Expression
)

type Parser struct {
    l           *lexer.Lexer
    curToken    token.Token
    peekToken   token.Token
    errors      []string

    prefixParseFns   map[token.TokenType]prefixParseFn
    infixParseFns    map[token.TokenType]infixParseFn
}

// precedence order
const (
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // +
    PRODUCT     // *
    PREFIX      // -X or !X
    CALL        // myFunction(X)
)

var precedencesMap = map[token.TokenType]int {
    token.EQ:       EQUALS,
    token.UNEQ:     EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}

func (p *Parser) peekPrecedence() int {
    if p, ok := precedencesMap[p.peekToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (p *Parser) curPrecedence() int {
    if p, ok := precedencesMap[p.curToken.Type]; ok {
        return p
    }
    return LOWEST
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l:      l,
        errors: []string{},
    }

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    p.registerPrefix(token.TRUE, p.parseBoolean)
    p.registerPrefix(token.FALSE, p.parseBoolean)
    p.registerPrefix(token.LPAREN, p.parseGroupedExpession)
    p.registerPrefix(token.IF, p.parseIfExpression)
    p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.UNEQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.LPAREN, p.parseCallExpression)

    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function for %s found", t)
    p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]

    if prefix == nil {
        p.noPrefixParseFnError(p.curToken.Type)
        return nil
    }

    left := prefix()
    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return left
        }
        p.nextToken()

        left = infix(left)
    }

    return left
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    e := &ast.PrefixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
    }

    p.nextToken()
    e.Right = p.parseExpression(PREFIX)

    return e
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    e := &ast.InfixExpression{
        Token:      p.curToken,
        Operator:   p.curToken.Literal,
        Left:       left,
    }

    precedence := p.curPrecedence()
    p.nextToken()
    e.Right = p.parseExpression(precedence)

    return e
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    s := &ast.ExpressionStatement{Token: p.curToken}

    s.Expression = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return s
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    s := &ast.ReturnStatement{Token: p.curToken}

    p.nextToken()

    s.ReturnValue = p.parseExpression(LOWEST)

    if !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return s
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    s := &ast.LetStatement{Token: p.curToken}

    if !p.expectPeek(token.IDENT) {
        return nil
    }

    s.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeek(token.ASSIGN) {
        return nil
    }
    p.nextToken()

    s.Value = p.parseExpression(LOWEST)

    for !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return s
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    l := &ast.IntegerLiteral{Token: p.curToken}

    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("could not parse %q as an int", p.curToken.Literal)
        p.errors = append(p.errors, msg)

        return nil
    }

    l.Value = value
    return l
}

func (p *Parser) parseBoolean() ast.Expression {
    return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpession() ast.Expression {
    p.nextToken()

    expression := p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAREN) {
        return nil
    }
    return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
    expression := &ast.IfExpression{Token: p.curToken}

    if !p.expectPeek(token.LPAREN) {
        return nil
    }

    p.nextToken()

    expression.Condition = p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    if !p.expectPeek(token.LBRACE) {
        return nil
    }

    expression.Consequence = p.parseBlockStatement()

    if p.peekTokenIs(token.ELSE) {
        p.nextToken()

        if !p.expectPeek(token.LBRACE) {
            return nil
        }

        expression.Alternative = p.parseBlockStatement()
    }

    return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
    fnLiteral := &ast.FunctionLiteral{Token: p.curToken}

    if !p.expectPeek(token.LPAREN) {
        return nil
    }

    fnLiteral.Parameters = p.parseFunctionParameters()

    if !p.expectPeek(token.LBRACE) {
        return nil
    }

    fnLiteral.Body = p.parseBlockStatement()

    return fnLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
    identifiers := []*ast.Identifier{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return identifiers
    }

    p.nextToken()

    identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    identifiers = append(identifiers, identifier)

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()

        identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
        identifiers = append(identifiers, identifier)
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    call := &ast.CallExpression{Token: p.curToken, Function: function}
    call.Arguments = p.parseCallArguments()

    return call
}

func (p *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return args
    }

    p.nextToken()
    args = append(args, p.parseExpression(LOWEST))

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()
        args = append(args, p.parseExpression(LOWEST))
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return args
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Token: p.curToken}
    block.Statements = []ast.Statement{}

    p.nextToken()

    for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
        statement := p.parseStatement()

        if statement != nil {
            block.Statements = append(block.Statements, statement)
        }

        p.nextToken()
    }
    return block
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
    msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        s := p.parseStatement()
        if s != nil {
            program.Statements = append(program.Statements, s)
        }
        p.nextToken()
    }

    return program
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

