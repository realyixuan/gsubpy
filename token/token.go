package token

type tokenType string

const (
    EOF         = "EOF"
    ILLEGAL     = "ILLEGAL"
    LINEFEED    = "LINEFEED"

    IF          = "if"
    WHILE       = "while"
    DEF         = "def"

    ASSIGN      = "ASSIGN"
    IDENTIFIER  = "IDENTIFIER"

    NUMBER      = "NUMBER"

    PLUS        = "PLUS"
    MINUS       = "MINUS"
    MUL         = "MUL"
    DIVIDE      = "DIVIDE"

    GT          = ">"
    LT          = "<"

    LPAREN      = "("
    RPAREN      = ")"
    COLON       = ":"
    COMMA       = ","
)

var Keywords = map[string]tokenType {
    "if":       IF,
    "while":    WHILE,
    "def":      DEF,
}


type Token struct {
    TokenType   tokenType
    Literals    string
}

