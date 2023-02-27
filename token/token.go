package token

type tokenType string

const (
    EOF         = "EOF"
    ILLEGAL     = "ILLEGAL"
    LINEFEED    = "LINEFEED"

    IF          = "if"
    WHILE       = "while"
    COLON       = ":"

    ASSIGN      = "ASSIGN"
    IDENTIFIER  = "IDENTIFIER"

    NUMBER      = "NUMBER"

    PLUS        = "PLUS"
    MINUS       = "MINUS"
    MUL         = "MUL"
    DIVIDE      = "DIVIDE"

    GT          = ">"
    LT          = "<"
)

var Keywords = map[string]tokenType {
    "if":       IF,
    "while":    WHILE,
}


type Token struct {
    TokenType   tokenType
    Literals    string
}

