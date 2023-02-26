package token

type tokenType string

const (
    EOF         = "EOF"
    ILLEGAL     = "ILLEGAL"
    LINEFEED    = "LINEFEED"

    ASSIGN      = "ASSIGN"
    IDENTIFIER  = "IDENTIFIER"

    NUMBER      = "NUMBER"

    PLUS        = "PLUS"
    MUL         = "MUL"
)


type Token struct {
    TokenType   tokenType
    Literals    string
}

