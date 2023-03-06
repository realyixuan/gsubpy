package token

type tokenType string

const (
    EOF         = "EOF"
    ILLEGAL     = "ILLEGAL"
    LINEFEED    = "LINEFEED"

    IF          = "if"
    WHILE       = "while"
    DEF         = "def"
    RETURN      = "return"

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
    "return":   RETURN,
}


type Token struct {
    Type   tokenType
    Literals    string
}

