package ast

import "github.com/UsamaHameed/monkey-interpreter/token"

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        // get the literal defined in the root node
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

type LetStatement struct {
    Token   token.Token
    Name    *Identifier
    Value   Expression
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

type Identifier struct {
    Token token.Token
    Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}

type ReturnStatement struct {
    Token   token.Token
    ReturnValue   Expression
}
func (ls *ReturnStatement) statementNode() {}
func (ls *ReturnStatement) TokenLiteral() string {
    return ls.Token.Literal
}
