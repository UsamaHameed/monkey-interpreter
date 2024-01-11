package parser

import (
	"fmt"
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integer, ok := il.(*ast.IntegerLiteral)

    if !ok {
        t.Errorf("il not an *ast.IntegerLiteral. got=%T", il)
        return false
    }

    if integer.Value != value {
        t.Errorf("integ.Value not %d. got=%d", value, integer.Value)
        return false
    }

    if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
        t.Errorf("integ.TokenLiteral not %d. got=%s", value,
            integer.TokenLiteral())
        return false
    }

    return true
}

func testBooleanLiteral(t *testing.T, il ast.Expression, value bool) bool {
    boolean, ok := il.(*ast.Boolean)

    if !ok {
        t.Errorf("expression not an *ast.Boolean. got=%T", il)
        return false
    }

    if boolean.Value != value {
        t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
        return false
    }

    if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
        t.Errorf("boolean.TokenLiteral not %t. got=%s", value,
            boolean.TokenLiteral())
        return false
    }

    return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
    identifier, ok := exp.(*ast.Identifier)

    if !ok {
        t.Errorf("expression is not an *ast.Identifier, got=%T", exp)
        return false
    }

    if identifier.Value != value {
        t.Errorf("identifier.Value not %s, got=%s", value, identifier.Value)
        return false
    }

    if identifier.TokenLiteral() != value {
        t.Errorf("identifier.TokenLiteral not %s, got=%s", value, identifier.TokenLiteral())
        return false
    }

    return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
    switch v := expected.(type) {
    case int:
        return testIntegerLiteral(t, exp, int64(v))
    case int64:
        return testIntegerLiteral(t, exp, v)
    case bool:
        return testBooleanLiteral(t, exp, v)
    case string:
        return testIdentifier(t, exp, v)
    }

    t.Errorf("type of expression not handled, got=%T", exp)
    return false
}

func testInfixExpression(
    t *testing.T,
    exp ast.Expression,
    left interface{},
    operator string,
    right interface{},
) bool {
    expression, ok := exp.(*ast.InfixExpression)

    if !ok {
        t.Errorf("expression is not ast.InfixExpression, got=%T(%s)", exp, exp)
        return false
    }

    if !testLiteralExpression(t, expression.Left, left) {
        return false
    }

    if expression.Operator != operator {
        t.Errorf("expression.Operator is not '%s', got=%q", operator, expression.Operator)
        return false
    }

    if !testLiteralExpression(t, expression.Right, right) {
        return false
    }

    return true
}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input      string
        leftValue  interface{}
        operator   string
        rightValue interface{}
    }{
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
        {"true == true", true, "==", true},
        {"true != false", true, "!=", false},
        {"false == false", false, "==", false},
    }

    for _, expected := range infixTests {
        l := lexer.New(expected.input)
        p := New(l)
        program := p.ParseProgram()
        checkParseErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                1, len(program.Statements))
        }

        s, ok := program.Statements[0].(*ast.ExpressionStatement)


        if !testInfixExpression(
            t,
            s.Expression,
            expected.leftValue,
            expected.operator,
            expected.rightValue,
        ) {
            return
        }

        if !ok {
            t.Fatalf("the root node aka program.Statements[0] is not an ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        if !ok {
            t.Fatalf("expression is not ast.InfixExpression. got=%T", s.Expression)
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {
            "-a * b",
            "((-a) * b)",
        },
        {
            "!-a",
            "(!(-a))",
        },
        {
            "a + b + c",
            "((a + b) + c)",
        },
        {
            "a + b - c",
            "((a + b) - c)",
        },
        {
            "a * b * c",
            "((a * b) * c)",
        },
        {
            "a * b / c",
            "((a * b) / c)",
        },
        {
            "a + b / c",
            "(a + (b / c))",
        },
        {
            "a + b * c + d / e - f",
            "(((a + (b * c)) + (d / e)) - f)",
        },
        {
            "3 + 4; -5 * 5",
            "(3 + 4)((-5) * 5)",
        },
        {
            "5 > 4 == 3 < 4",
            "((5 > 4) == (3 < 4))",
        },
        {
            "5 < 4 != 3 > 4",
            "((5 < 4) != (3 > 4))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "true",
            "true",
        },
        {
            "false",
            "false",
        },
        {
            "3 > 5 == false",
            "((3 > 5) == false)",
        },
        {
            "3 < 5 == true",
            "((3 < 5) == true)",
        },
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        p := New(l)
        program := p.ParseProgram()
        checkParseErrors(t, p)

        actual := program.String()
        if actual != test.expected {
            t.Errorf("expected=%q, got=%q", test.expected, actual)
        }
    }
}

func TestBooleanExpression(t *testing.T) {
    tests := []struct {
        input       string
        expected    bool
    } {
        {"true", true},
        {"false", false},
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()

        checkParseErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not have 1 statements, got=%d",
            len(program.Statements))
        }

        s, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("the root node aka program.Statements[0] is not an ExpressionStatement, got=%T",
            program.Statements[0])
        }

        boolean, ok := s.Expression.(*ast.Boolean)

        if !ok {
            t.Fatalf("expression is not a boolean, got=%T", s.Expression)
        }

        if boolean.Value != test.expected {
            t.Fatalf("boolean value not %t, got=%t", test.expected, boolean.Value)
        }
    }
}
