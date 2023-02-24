package lexer

import (
    "gsubpy/token"
)

import (
    "testing"
)

func TestLexer(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `val = 1`,
            []token.Token{
                token.Token{TokenType: token.IDENTIFIER, Literals: "val"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "1"},
                token.Token{TokenType: token.EOF},
            },
        },
        {
            `val = 1 + 1`,
            []token.Token{
                token.Token{TokenType: token.IDENTIFIER, Literals: "val"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "1"},
                token.Token{TokenType: token.PLUS, Literals: "+"},
                token.Token{TokenType: token.NUMBER, Literals: "1"},
                token.Token{TokenType: token.EOF},
            },
        },
        {
            `val = 10 + 20 * 10`,
            []token.Token{
                token.Token{TokenType: token.IDENTIFIER, Literals: "val"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "10"},
                token.Token{TokenType: token.PLUS, Literals: "+"},
                token.Token{TokenType: token.NUMBER, Literals: "20"},
                token.Token{TokenType: token.MUL, Literals: "*"},
                token.Token{TokenType: token.NUMBER, Literals: "10"},
                token.Token{TokenType: token.EOF},
            },
        },
    }

    for _, testCase := range testCases {
        l := New(testCase.input)
        for _, tk := range testCase.expectedTokens {
            if tk != l.CurToken {
                t.Errorf("expected token %s, got token %s", tk, l.CurToken)
            }
            l.ReadNextToken()
        }
    }
}

