package lexer

import (
    "gsubpy/token"
    "gsubpy/object"
)

/*
the lexer have following aiblities: 
    - view current token of lexer
    - read next token
    - peek next token
    
*/

type Lexer struct {
    input       string
    idx         int
    ch          byte
    Indents     string
    indentReady bool
    CurToken    token.Token
}

func New(input string) *Lexer {
    l := &Lexer{input: input, indentReady: true}
    l.readChar()
    l.ReadNextToken()
    return l
}

func (l *Lexer) ReadNextToken() {
    l.skipWhitespace()

    switch l.ch {
    case '=':
        l.CurToken = token.Token{Type: token.ASSIGN, Literals: string(l.ch)}
        l.readChar()
    case '+':
        l.CurToken = token.Token{Type: token.PLUS, Literals: string(l.ch)}
        l.readChar()
    case '-':
        l.CurToken = token.Token{Type: token.MINUS, Literals: string(l.ch)}
        l.readChar()
    case '*':
        l.CurToken = token.Token{Type: token.MUL, Literals: string(l.ch)}
        l.readChar()
    case '/':
        l.CurToken = token.Token{Type: token.DIVIDE, Literals: string(l.ch)}
        l.readChar()
    case '>':
        l.CurToken = token.Token{Type: token.GT, Literals: string(l.ch)}
        l.readChar()
    case '<':
        l.CurToken = token.Token{Type: token.LT, Literals: string(l.ch)}
        l.readChar()
    case '(':
        l.CurToken = token.Token{Type: token.LPAREN, Literals: string(l.ch)}
        l.readChar()
    case ')':
        l.CurToken = token.Token{Type: token.RPAREN, Literals: string(l.ch)}
        l.readChar()
    case '[':
        l.CurToken = token.Token{Type: token.LBRACKET, Literals: string(l.ch)}
        l.readChar()
    case ']':
        l.CurToken = token.Token{Type: token.RBRACKET, Literals: string(l.ch)}
        l.readChar()
    case '{':
        l.CurToken = token.Token{Type: token.LBRACE, Literals: string(l.ch)}
        l.readChar()
    case '}':
        l.CurToken = token.Token{Type: token.RBRACE, Literals: string(l.ch)}
        l.readChar()
    case ',':
        l.CurToken = token.Token{Type: token.COMMA, Literals: string(l.ch)}
        l.readChar()
    case ':':
        l.CurToken = token.Token{Type: token.COLON, Literals: string(l.ch)}
        l.readChar()
    case '.':
        l.CurToken = token.Token{Type: token.DOT, Literals: string(l.ch)}
        l.readChar()
    case '"', '\'':
        l.CurToken = token.Token{Type: token.STRING}
        l.CurToken.Literals = l.readString()
        l.readChar()
    case '\n':
        l.CurToken = token.Token{Type: token.LINEFEED, Literals: string(l.ch)}
        l.readChar()
        l.indentReady = true
    case '\x03':
        l.CurToken = token.Token{Type: token.EOF}
    default:
        if isDigit(l.ch) {
            num := l.readNumber()
            l.CurToken = token.Token{Type: token.INTEGER, Literals: num}
        } else if isLetter(l.ch) {
            identifier := l.readLetter()
            if tokType, ok := token.Keywords[identifier]; ok {
                l.CurToken = token.Token{Type: tokType, Literals: identifier}
            } else {
                l.CurToken = token.Token{Type: token.IDENTIFIER, Literals: identifier}
            }
        } else {
            panic(&object.ExceptionInst{Msg: "syntaxException: syntax error"})

            l.CurToken = token.Token{Type: token.ILLEGAL}
            l.readChar()
        }
    }
}

func (l *Lexer) PeekNextToken() token.Token {
    // A trick
    lc := *l
    lc.ReadNextToken()
    return lc.CurToken
}

func (l *Lexer) skipWhitespace() {
    var indentStr string
    for l.ch == ' ' || l.ch == '\t' {
        indentStr += string(l.ch)
        l.readChar()
    }

    if l.indentReady {
        l.Indents = indentStr
        l.indentReady = false
    }
}

func (l *Lexer) readString() string {
    var result string
    stringMark := l.ch
    l.readChar()
    for l.ch != stringMark && l.ch != '\x03' {
        result += string(l.ch)
        l.readChar()
    }

    if l.ch != stringMark {
        panic("string have wrong format")
    }

    return result
}

func (l *Lexer) readNumber() string {
    res := ""
    for isDigit(l.ch) {
        res += string(l.ch)
        l.readChar()
    }
    return res
}

func (l *Lexer) readLetter() string {
    res := ""
    for isLetter(l.ch) || isDigit(l.ch) {
        res += string(l.ch)
        l.readChar()
    }
    return res
}

func (l *Lexer) readChar() {
    if l.idx < len(l.input) {
        l.ch = l.input[l.idx]
    } else {
        l.ch = '\x03'   // end of text, special byte
    }
    l.idx += 1
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}


