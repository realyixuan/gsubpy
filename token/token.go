package token

type tokenType string

const (
    EOF         = "EOF"
    ASSIGN      = "="
    PLUS        = "+"
    IDENTIFIER  = "IDENTIFIER"
    NUMBER      = "NUMBER"
)


type Token struct {
    TokenType   tokenType
    Literals    string
}

