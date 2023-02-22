package token

type tokenType string

const (
    EOF         = "EOF"
    ASSIGN      = "ASSIGN"
    PLUS        = "PLUS"
    IDENTIFIER  = "IDENTIFIER"
    NUMBER      = "NUMBER"
)


type Token struct {
    TokenType   tokenType
    Literals    string
}

