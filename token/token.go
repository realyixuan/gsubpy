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
    FOR         = "for"
    IN          = "in"
    NIN         = "not in"
    IS          = "is"
    ISN         = "is not"
    ASSERT      = "assert"

    ASSIGN      = "="
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
    EQ          = "=="

    AND         = "and"
    OR          = "or"
    NOT         = "not"

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
    "and":      AND,
    "or":       OR,
    "not":      NOT,
    "for":      FOR,
    "in":       IN,
    "not in":   NIN,
    "is":       IS,
    "is not":   ISN,
    "assert":   ASSERT,
}


type Token struct {
    Type   TokenType
    Literals    string
}

