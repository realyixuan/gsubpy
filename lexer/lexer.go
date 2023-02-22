package lexer

import (
    "gsubpy/token"
)

type Lexer struct {
    Input   string
    idx     int
}

func (l *Lexer) NextToken() token.Token {
    //  get tokens
    l.skipWhitespace()

    if l.isExhausted() {
        return token.Token{TokenType: token.EOF}
    }

    switch ch := l.Input[l.idx]; ch {
    case '=':
        tk := token.Token{TokenType: token.ASSIGN, Literals: string(ch)}
        l.idx += 1
        return tk
    default:
        if isDigit(ch) {
            num := l.readNumber()
            tk := token.Token{TokenType: token.NUMBER, Literals: num}
            l.idx += 1
            return tk
        } else if isLetter(ch) {
            identifier := l.readLetter()
            tk := token.Token{TokenType: token.IDENTIFIER, Literals: identifier}
            l.idx += 1
            return tk
        }
    }

    return token.Token{TokenType: token.EOF}

}

func (l *Lexer) peekNextToken() {
    //
}

func (l *Lexer) skipWhitespace() {
    for !l.isExhausted() && (l.Input[l.idx] == ' ' || l.Input[l.idx] == '\t' || l.Input[l.idx] == '\r') {
        l.idx += 1
    }
}

func (l *Lexer) isExhausted() bool {
    return l.idx >= len(l.Input)
}

func (l *Lexer) readNumber() string {
    res := ""
    for !l.isExhausted() && isDigit(l.Input[l.idx]) {
        res += string(l.Input[l.idx])
        l.idx += 1
    }
    return res
}

func (l *Lexer) readLetter() string {
    res := ""
    for !l.isExhausted() && isLetter(l.Input[l.idx]) {
        res += string(l.Input[l.idx])
        l.idx += 1
    }
    return res
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}


