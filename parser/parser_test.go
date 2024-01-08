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
    return a + b;
    return ++s;
`
    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)

    if program == nil {
        t.Fatalf("ParseProgram returned nil")
    }

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements length is not 3, got %d", len(program.Statements))
    }

    for _, s := range program.Statements {
        statement, ok := s.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("statement not *ast.ReturnStatement, got=%T", s)
            continue
        }
        if statement.TokenLiteral() != "return" {
            t.Errorf("return statement TokenLiteral not 'return', got %q", statement.TokenLiteral())
        }
    }
}
