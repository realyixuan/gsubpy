package token

type TokenType string

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
    CLASS       = "class"

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
    LBRACE      = "{"
    RBRACE      = "}"
    COLON       = ":"
    COMMA       = ","
)

var Keywords = map[string]TokenType {
    "if":       IF,
    "elif":     ELIF,
    "else":     ELSE,
    "while":    WHILE,
    "def":      DEF,
    "return":   RETURN,
    "class":    CLASS,
}


type Token struct {
    Type   TokenType
    Literals    string
}

