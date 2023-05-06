package lexer

import (
    "fmt"

    "github.com/realyixuan/gsubpy/token"
    "github.com/realyixuan/gsubpy/evaluator"
)

type Lexer struct {
    input       string
    idx         int
    LineNum     int
    Line        string
    ch          byte
    Indents     string
    indentReady bool
    CurToken    token.Token
}

func New(input string) *Lexer {
    l := &Lexer{input: input, indentReady: true}
    l.readLine()
    l.readChar()
    l.ReadNextToken()
    return l
}

func (l *Lexer) ReadNextToken() {
    l.readNextToken()

    if _, ok := token.Keywords[l.CurToken.Literals]; ok {
        cl := *l
        cl.readNextToken()
        if _, ok := token.Keywords[cl.CurToken.Literals]; ok {
            letters := l.CurToken.Literals + " " + cl.CurToken.Literals
            if tokType, ok := token.Keywords[letters]; ok {
                l.readNextToken()
                l.CurToken = token.Token{Type: tokType, Literals: letters}
            }
        }
    }


}

func (l *Lexer) readNextToken() {
    l.skipWhitespace()
    l.skipoverComment()

    switch l.ch {
    case '=':
        l.readChar()
        if l.ch == '=' {
            l.CurToken = token.Token{Type: token.EQ, Literals: token.EQ}
            l.readChar()
        } else {
            l.CurToken = token.Token{Type: token.ASSIGN, Literals: token.ASSIGN}
        }
    case '!':
        l.readChar()
        if l.ch == '=' {
            l.CurToken = token.Token{Type: token.NEQ, Literals: token.NEQ}
            l.readChar()
        } else {
            panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: invalid syntax", l.LineNum, l.Line)))
        }
    case '+':
        l.readChar()
        if l.ch == '=' {
            l.CurToken = token.Token{Type: token.PLUSASSIGN, Literals: token.PLUSASSIGN}
            l.readChar()
        } else {
            l.CurToken = token.Token{Type: token.PLUS, Literals: token.PLUS}
        }
    case '-':
        l.readChar()
        if l.ch == '=' {
            l.CurToken = token.Token{Type: token.MINUSASSIGN, Literals: token.MINUSASSIGN}
            l.readChar()
        } else {
            l.CurToken = token.Token{Type: token.MINUS, Literals: "-"}
        }
    case '*':
        l.readChar()
        if l.ch == '=' {
            l.CurToken = token.Token{Type: token.MULASSIGN, Literals: token.MULASSIGN}
            l.readChar()
        } else {
            l.CurToken = token.Token{Type: token.MUL, Literals: "*"}
        }
    case '/':
        l.readChar()
        if l.ch == '=' {
            l.CurToken = token.Token{Type: token.DIVIDEASSIGN, Literals: token.DIVIDEASSIGN}
            l.readChar()
        } else {
            l.CurToken = token.Token{Type: token.DIVIDE, Literals: token.DIVIDE}
        }
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
    case '\\':
        l.readChar()
        if l.ch != '\n' {
            panic(evaluator.Error(fmt.Sprintf(
                "line %v\n\t%s\nSyntaxError: unexpected character after line continuation character",
                l.LineNum, l.Line)))
        }
        l.readChar()
        l.LineNum++
        l.readNextToken()
    case '\n':
        l.CurToken = token.Token{Type: token.LINEFEED, Literals: string(l.ch)}
        l.readLine()
        l.readChar()
        l.indentReady = true
    case '\x03':
        l.CurToken = token.Token{Type: token.EOF}
    default:
        if isDigit(l.ch) {
            num := l.readNumber()
            l.CurToken = token.Token{Type: token.INTEGER, Literals: num}
        } else if isLetter(l.ch) {
            identifier := l.readLetters()
            if tokType, ok := token.Keywords[identifier]; ok {
                l.CurToken = token.Token{Type: tokType, Literals: identifier}
            } else {
                l.CurToken = token.Token{Type: token.IDENTIFIER, Literals: identifier}
            }
        } else {
            panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: invalid syntax", l.LineNum, l.Line)))

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
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: invalid syntax, string have wrong format", l.LineNum, l.Line)))
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

func (l *Lexer) readLetters() string {
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

func (l *Lexer)readLine() {
    l.Line = ""
    for idx := l.idx; idx < len(l.input) && l.input[idx] != '\n'; idx++ {
        if l.input[idx] == '\\' {
            if idx+1 < len(l.input) && l.input[idx+1] == '\n' {
                l.Line += string(l.input[idx:idx+2])
                idx++
            } else {
                l.Line += string(l.input[idx])
            }
        } else {
            l.Line += string(l.input[idx])
        }
    }
    l.LineNum++
}

func (l *Lexer) skipoverComment() {
    if l.ch == '#' {
        for l.ch != '\n' && l.ch != '\x03' {
            l.readChar()
        }
    }
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}


