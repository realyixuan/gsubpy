package lexer

import (
    "gsubpy/token"
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
        l.CurToken = token.Token{TokenType: token.ASSIGN, Literals: string(l.ch)}
        l.readChar()
    case '+':
        l.CurToken = token.Token{TokenType: token.PLUS, Literals: string(l.ch)}
        l.readChar()
    case '-':
        l.CurToken = token.Token{TokenType: token.MINUS, Literals: string(l.ch)}
        l.readChar()
    case '*':
        l.CurToken = token.Token{TokenType: token.MUL, Literals: string(l.ch)}
        l.readChar()
    case '/':
        l.CurToken = token.Token{TokenType: token.DIVIDE, Literals: string(l.ch)}
        l.readChar()
    case '>':
        l.CurToken = token.Token{TokenType: token.GT, Literals: string(l.ch)}
        l.readChar()
    case '<':
        l.CurToken = token.Token{TokenType: token.LT, Literals: string(l.ch)}
        l.readChar()
    case '(':
        l.CurToken = token.Token{TokenType: token.LPAREN, Literals: string(l.ch)}
        l.readChar()
    case ')':
        l.CurToken = token.Token{TokenType: token.RPAREN, Literals: string(l.ch)}
        l.readChar()
    case ',':
        l.CurToken = token.Token{TokenType: token.COMMA, Literals: string(l.ch)}
        l.readChar()
    case ':':
        l.CurToken = token.Token{TokenType: token.COLON, Literals: string(l.ch)}
        l.readChar()
    case '\n':
        l.CurToken = token.Token{TokenType: token.LINEFEED, Literals: string(l.ch)}
        l.readChar()
        l.indentReady = true
    case '\x03':
        l.CurToken = token.Token{TokenType: token.EOF}
    default:
        if isDigit(l.ch) {
            num := l.readNumber()
            l.CurToken = token.Token{TokenType: token.NUMBER, Literals: num}
        } else if isLetter(l.ch) {
            identifier := l.readLetter()
            if tokType, ok := token.Keywords[identifier]; ok {
                l.CurToken = token.Token{TokenType: tokType, Literals: identifier}
            } else {
                l.CurToken = token.Token{TokenType: token.IDENTIFIER, Literals: identifier}
            }
        } else {
            l.CurToken = token.Token{TokenType: token.ILLEGAL}
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
    for isLetter(l.ch) {
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
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}


