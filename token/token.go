package token

type TokenType string

type Token struct {
    Type        TokenType
    Literal     string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"

    IDENT = "IDENT"

    ASSIGN = "="
    PLUS = "+"
    MINUS = "-"

    INT = "INT"
)

var keywords = map[string]TokenType{
    "...": ILLEGAL,
}

func LookupIdent(ident string) TokenType {
    if toktype, ok := keywords[ident]; ok {
        return toktype
    }
    return IDENT
}

