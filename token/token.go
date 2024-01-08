package token

type TokenType string

type Token struct {
    Type    TokenType
    Literal string
}

const (
    ILLEGAL     = "ILLEGAL"
    EOF         = "EOF"
    IDENT       = "IDENT"
    INT         = "INT"
    COMMA       = ","
    ASSIGN      = "="
    PLUS        = "+"
    MINUS       = "-"
    BANG        = "!"
    ASTERISK    = "*"
    SLASH       = "/"
    LT          = "<"
    GT          = ">"
    SEMICOLON   = ";"
    LPAREN      = "("
    RPAREN      = ")"
    LBRACE      = "{"
    RBRACE      = "}"
    FUNCTION    = "FUNCTION"
    LET         = "LET"
    IF          = "IF"
    TRUE        = "TRUE"
    FALSE       = "FALSE"
    ELSE        = "ELSE"
    RETURN      = "RETURN"
    EQ          = "=="
    UNEQ        = "!="
)

var keywords = map[string]TokenType{
    "fn":       FUNCTION,
    "let":      LET,
    "if":       IF,
    "true":     TRUE,
    "false":    FALSE,
    "else":     ELSE,
    "return":   RETURN,
}

func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }

    return IDENT
}
