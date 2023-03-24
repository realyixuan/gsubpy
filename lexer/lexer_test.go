package lexer

import (
    "gsubpy/token"
)

import (
    "testing"
)

func TestIdentifier(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `_abc12`,
            []token.Token{
                token.Token{Type: token.IDENTIFIER, Literals: "_abc12"},
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

func TestAssignStatement(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `val = 10 + 20 * 10 / 2 - 50`,
            []token.Token{
                token.Token{Type: token.IDENTIFIER, Literals: "val"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.INTEGER, Literals: "10"},
                token.Token{Type: token.PLUS, Literals: "+"},
                token.Token{Type: token.INTEGER, Literals: "20"},
                token.Token{Type: token.MUL, Literals: "*"},
                token.Token{Type: token.INTEGER, Literals: "10"},
                token.Token{Type: token.DIVIDE, Literals: "/"},
                token.Token{Type: token.INTEGER, Literals: "2"},
                token.Token{Type: token.MINUS, Literals: "-"},
                token.Token{Type: token.INTEGER, Literals: "50"},
                token.Token{Type: token.EOF},
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
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.INTEGER, Literals: "1"},
                token.Token{Type: token.LINEFEED, Literals: "\n"},
                token.Token{Type: token.IDENTIFIER, Literals: "b"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.INTEGER, Literals: "2"},
                token.Token{Type: token.LINEFEED, Literals: "\n"},
                token.Token{Type: token.EOF},
            },
        },
        {
            "a = 1\n" +
            "b = a + 2\n",
            []token.Token{
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.INTEGER, Literals: "1"},
                token.Token{Type: token.LINEFEED, Literals: "\n"},
                token.Token{Type: token.IDENTIFIER, Literals: "b"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.PLUS, Literals: "+"},
                token.Token{Type: token.INTEGER, Literals: "2"},
                token.Token{Type: token.LINEFEED, Literals: "\n"},
                token.Token{Type: token.EOF},
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
                token.Token{Type: token.IF, Literals: "if"},
                token.Token{Type: token.INTEGER, Literals: "10"},
                token.Token{Type: token.GT, Literals: ">"},
                token.Token{Type: token.INTEGER, Literals: "5"},
                token.Token{Type: token.COLON, Literals: ":"},
                token.Token{Type: token.LINEFEED, Literals: "\n"},
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.INTEGER, Literals: "1"},
                token.Token{Type: token.LINEFEED, Literals: "\n"},
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
                token.Token{Type: token.WHILE, Literals: "while"},
                token.Token{Type: token.INTEGER, Literals: "10"},
                token.Token{Type: token.GT, Literals: ">"},
                token.Token{Type: token.INTEGER, Literals: "5"},
                token.Token{Type: token.COLON, Literals: ":"},
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
                token.Token{Type: token.DEF, Literals: "def"},
                token.Token{Type: token.IDENTIFIER, Literals: "foo"},
                token.Token{Type: token.LPAREN, Literals: "("},
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.COMMA, Literals: ","},
                token.Token{Type: token.IDENTIFIER, Literals: "b"},
                token.Token{Type: token.RPAREN, Literals: ")"},
                token.Token{Type: token.COLON, Literals: ":"},
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

func TestReturnStatement(t *testing.T) {
    testCases := []struct{
        input           string
        expectedTokens  []token.Token
    } {
        {
            "return a + b",
            []token.Token{
                token.Token{Type: token.RETURN, Literals: "return"},
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.PLUS, Literals: "+"},
                token.Token{Type: token.IDENTIFIER, Literals: "b"},
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

func TestStringStatement(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `val = "abc,de"`,
            []token.Token{
                token.Token{Type: token.IDENTIFIER, Literals: "val"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.STRING, Literals: "abc,de"},
                token.Token{Type: token.EOF},
            },
        },
        {
            `val = 'abc,de'`,
            []token.Token{
                token.Token{Type: token.IDENTIFIER, Literals: "val"},
                token.Token{Type: token.ASSIGN, Literals: "="},
                token.Token{Type: token.STRING, Literals: "abc,de"},
                token.Token{Type: token.EOF},
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

func TestList(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `["abc", a, 123, ""]`,
            []token.Token{
                token.Token{Type: token.LBRACKET, Literals: "["},
                token.Token{Type: token.STRING, Literals: "abc"},
                token.Token{Type: token.COMMA, Literals: ","},
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.COMMA, Literals: ","},
                token.Token{Type: token.INTEGER, Literals: "123"},
                token.Token{Type: token.COMMA, Literals: ","},
                token.Token{Type: token.STRING, Literals: ""},
                token.Token{Type: token.RBRACKET, Literals: "]"},
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

func TestDict(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            `{'a': 1, '10': 10, var: "abc"}`,
            []token.Token{
                token.Token{Type: token.LBRACE, Literals: "{"},
                token.Token{Type: token.STRING, Literals: "a"},
                token.Token{Type: token.COLON, Literals: ":"},
                token.Token{Type: token.INTEGER, Literals: "1"},
                token.Token{Type: token.COMMA, Literals: ","},
                token.Token{Type: token.STRING, Literals: "10"},
                token.Token{Type: token.COLON, Literals: ":"},
                token.Token{Type: token.INTEGER, Literals: "10"},
                token.Token{Type: token.COMMA, Literals: ","},
                token.Token{Type: token.IDENTIFIER, Literals: "var"},
                token.Token{Type: token.COLON, Literals: ":"},
                token.Token{Type: token.STRING, Literals: "abc"},
                token.Token{Type: token.RBRACE, Literals: "}"},
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

func TestClass(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            "class Foo:",
            []token.Token{
                token.Token{Type: token.CLASS, Literals: "class"},
                token.Token{Type: token.IDENTIFIER, Literals: "Foo"},
                token.Token{Type: token.COLON, Literals: ":"},
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

func TestDot(t *testing.T) {
    testCases := []struct {
        input   string
        expectedTokens []token.Token
    }{
        {
            "a.b.c",
            []token.Token{
                token.Token{Type: token.IDENTIFIER, Literals: "a"},
                token.Token{Type: token.DOT, Literals: "."},
                token.Token{Type: token.IDENTIFIER, Literals: "b"},
                token.Token{Type: token.DOT, Literals: "."},
                token.Token{Type: token.IDENTIFIER, Literals: "c"},
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

