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
    UNDERSCORE  = "_"

    INTEGER      = "INTEGER"
    STRING      = "STRING"

    PLUS        = "+"
    MINUS       = "-"
    MUL         = "*"
    DIVIDE      = "/"
    PLUSASSIGN  = "+="
    MINUSASSIGN  = "-="
    MULASSIGN  = "*="
    DIVIDEASSIGN  = "/="

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
    DOT         = "."
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

