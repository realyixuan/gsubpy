package token

type tokenType string

const (
    EOF         = "EOF"
    ILLEGAL     = "ILLEGAL"
    LINEFEED    = "LINEFEED"

    IF          = "if"
    ELIF        = "elif"
    ELSE        = "else"
    WHILE       = "while"
    DEF         = "def"
    RETURN      = "return"

    ASSIGN      = "ASSIGN"
    IDENTIFIER  = "IDENTIFIER"

    NUMBER      = "NUMBER"
    STRING      = "STRING"

    PLUS        = "PLUS"
    MINUS       = "MINUS"
    MUL         = "MUL"
    DIVIDE      = "DIVIDE"

    GT          = ">"
    LT          = "<"

    LPAREN      = "("
    RPAREN      = ")"
    LBRACKET    = "["
    RBRACKET    = "]"
    COLON       = ":"
    COMMA       = ","
)

var Keywords = map[string]tokenType {
    "if":       IF,
    "elif":     ELIF,
    "else":     ELSE,
    "while":    WHILE,
    "def":      DEF,
    "return":   RETURN,
}


type Token struct {
    Type   tokenType
    Literals    string
}

