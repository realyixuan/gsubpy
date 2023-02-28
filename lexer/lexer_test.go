package lexer

import (
    "gsubpy/token"
)

import (
    "testing"
)

func TestAssignStatement(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `val = 10 + 20 * 10 / 2 - 50`,
            []token.Token{
                token.Token{TokenType: token.IDENTIFIER, Literals: "val"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "10"},
                token.Token{TokenType: token.PLUS, Literals: "+"},
                token.Token{TokenType: token.NUMBER, Literals: "20"},
                token.Token{TokenType: token.MUL, Literals: "*"},
                token.Token{TokenType: token.NUMBER, Literals: "10"},
                token.Token{TokenType: token.DIVIDE, Literals: "/"},
                token.Token{TokenType: token.NUMBER, Literals: "2"},
                token.Token{TokenType: token.MINUS, Literals: "-"},
                token.Token{TokenType: token.NUMBER, Literals: "50"},
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

func TestMultiLineStatement(t *testing.T) {
    testCases := []struct{
        input           string
        expectedTokens  []token.Token
    } {
        {
            "a = 1\n" +
            "b = 2\n",
            []token.Token{
                token.Token{TokenType: token.IDENTIFIER, Literals: "a"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "1"},
                token.Token{TokenType: token.LINEFEED, Literals: "\n"},
                token.Token{TokenType: token.IDENTIFIER, Literals: "b"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "2"},
                token.Token{TokenType: token.LINEFEED, Literals: "\n"},
                token.Token{TokenType: token.EOF},
            },
        },
        {
            "a = 1\n" +
            "b = a + 2\n",
            []token.Token{
                token.Token{TokenType: token.IDENTIFIER, Literals: "a"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "1"},
                token.Token{TokenType: token.LINEFEED, Literals: "\n"},
                token.Token{TokenType: token.IDENTIFIER, Literals: "b"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.IDENTIFIER, Literals: "a"},
                token.Token{TokenType: token.PLUS, Literals: "+"},
                token.Token{TokenType: token.NUMBER, Literals: "2"},
                token.Token{TokenType: token.LINEFEED, Literals: "\n"},
                token.Token{TokenType: token.EOF},
            },
        },
    }


    for _, testCase := range testCases {
        l := New(testCase.input)
        for _, tk := range testCase.expectedTokens {
            if tk != l.CurToken {
                t.Errorf("expected %s, got %s", tk, l.CurToken)
            }
            l.ReadNextToken()
        }
    }

}

func TestIfStatement(t *testing.T) {
    testCases := []struct{
        input           string
        expectedTokens  []token.Token
    } {
        {
            "if 10 > 5:\n" +
            "    a = 1\n",
            []token.Token{
                token.Token{TokenType: token.IF, Literals: "if"},
                token.Token{TokenType: token.NUMBER, Literals: "10"},
                token.Token{TokenType: token.GT, Literals: ">"},
                token.Token{TokenType: token.NUMBER, Literals: "5"},
                token.Token{TokenType: token.COLON, Literals: ":"},
                token.Token{TokenType: token.LINEFEED, Literals: "\n"},
                token.Token{TokenType: token.IDENTIFIER, Literals: "a"},
                token.Token{TokenType: token.ASSIGN, Literals: "="},
                token.Token{TokenType: token.NUMBER, Literals: "1"},
                token.Token{TokenType: token.LINEFEED, Literals: "\n"},
            },
        },
    }


    for _, testCase := range testCases {
        l := New(testCase.input)
        for _, tk := range testCase.expectedTokens {
            if tk != l.CurToken {
                t.Errorf("expected %s, got %s", tk, l.CurToken)
            }
            l.ReadNextToken()
        }
    }

}

func TestWhileStatement(t *testing.T) {
    testCases := []struct{
        input           string
        expectedTokens  []token.Token
    } {
        {
            "while 10 > 5:",
            []token.Token{
                token.Token{TokenType: token.WHILE, Literals: "while"},
                token.Token{TokenType: token.NUMBER, Literals: "10"},
                token.Token{TokenType: token.GT, Literals: ">"},
                token.Token{TokenType: token.NUMBER, Literals: "5"},
                token.Token{TokenType: token.COLON, Literals: ":"},
            },
        },
    }


    for _, testCase := range testCases {
        l := New(testCase.input)
        for _, tk := range testCase.expectedTokens {
            if tk != l.CurToken {
                t.Errorf("expected %s, got %s", tk, l.CurToken)
            }
            l.ReadNextToken()
        }
    }

}

func TestDefStatement(t *testing.T) {
    testCases := []struct{
        input           string
        expectedTokens  []token.Token
    } {
        {
            "def foo(a, b):",
            []token.Token{
                token.Token{TokenType: token.DEF, Literals: "def"},
                token.Token{TokenType: token.IDENTIFIER, Literals: "foo"},
                token.Token{TokenType: token.LPAREN, Literals: "("},
                token.Token{TokenType: token.IDENTIFIER, Literals: "a"},
                token.Token{TokenType: token.COMMA, Literals: ","},
                token.Token{TokenType: token.IDENTIFIER, Literals: "b"},
                token.Token{TokenType: token.RPAREN, Literals: ")"},
                token.Token{TokenType: token.COLON, Literals: ":"},
            },
        },
    }


    for _, testCase := range testCases {
        l := New(testCase.input)
        for _, tk := range testCase.expectedTokens {
            if tk != l.CurToken {
                t.Errorf("expected %s, got %s", tk, l.CurToken)
            }
            l.ReadNextToken()
        }
    }

}

