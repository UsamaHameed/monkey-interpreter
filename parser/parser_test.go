package parser

import (
	"testing"

	"github.com/UsamaHameed/monkey-interpreter/ast"
	"github.com/UsamaHameed/monkey-interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
    input := `
let x = 5;
let y = 10;
let a = 1;
let foo = 838383;
`
    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)

    if program == nil {
        t.Fatalf("ParseProgram returned nil")
    }

    if len(program.Statements) != 4 {
        t.Fatalf("program.Statements length is not 4, got %d", len(program.Statements))
    }

    tests := []struct {
        expectedIdentifier string
    } {
        {"x"},
        {"y"},
        {"a"},
        {"foo"},
    }

    for i, expected := range tests {

        s := program.Statements[i]
        if !testLetStatement(t, s, expected.expectedIdentifier) {
            return
        }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, name string)  bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral did not return 'let'. got=%q", s.TokenLiteral())
        return false
    }

    // what is this???
    statement, ok := s.(*ast.LetStatement)

    if !ok {
        t.Errorf("s not *ast.LetStatement. got=%T", s)
    }

    if statement.Name.Value != name {
        t.Errorf("LetStatement.Name.Value not '%s'. got='%s'", name, statement.Name.Value)
        return false
    }

    if statement.Name.TokenLiteral() != name {
        t.Errorf("LetStatement.Name.TokenLiteral() not '%s'. got='%s'", name, statement.Name.TokenLiteral())
        return false
    }

    return true
}

func checkParseErrors(t *testing.T, p *Parser) {
    errors := p.Errors()
    if len(errors) == 0 {
        return
    }

    t.Errorf("%d parser errors", len(errors))

    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}


func TestReturnStatement(t *testing.T) {
    input := `
    return 1;
    return 10;
    return 9999999;
`
    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)

    if program == nil {
        t.Fatalf("ParseProgram returned nil")
    }

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements length is not 3, got %d",
        len(program.Statements))
    }

    for _, s := range program.Statements {
        statement, ok := s.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("statement not *ast.ReturnStatement, got=%T", s)
            continue
        }
        if statement.TokenLiteral() != "return" {
            t.Errorf("return statement TokenLiteral not 'return', got %q",
            statement.TokenLiteral())
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()

    checkParseErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program did not return enough statements, got=%d",
        len(program.Statements))
    }

    s, ok := program.Statements[0].(*ast.ExpressionStatement)

    if !ok {
        t.Fatalf("the root node aka program.Statements[0] is not an ast.ExpressionStatement, got=%T",
        program.Statements[0])
    }

    identifier, ok := s.Expression.(*ast.Identifier)

    if !ok {
        t.Fatalf("expression is not an *ast.Identifier, got=%T",
        program.Statements[0])
    }

    if identifier.Value != "foobar" {
        t.Errorf("identifier.Value is not %s, got=%s", "foobar",
        identifier.Value)
    }

    if identifier.TokenLiteral() != "foobar" {
        t.Errorf("identifier.TokenLiteral not %s, got %s", "foobar",
        identifier.TokenLiteral())
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program did not return enough statements. got=%d",
            len(program.Statements))
    }
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("the root node aka program.Statements[0] is not an ast.ExpressionStatement, got=%T",
            program.Statements[0])
    }

    literal, ok := stmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
    }
    if literal.Value != 5 {
        t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
    }
    if literal.TokenLiteral() != "5" {
        t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
            literal.TokenLiteral())
    }
}
