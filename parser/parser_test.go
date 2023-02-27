package parser

import (
    "testing"

    "gsubpy/lexer"
)

func TestIndentParsing(t *testing.T) {
    input := "" +
    "    a =       1\n" + 
    "a = 1\n" + 
    "    a =     1\n" + 
    "        a      = 1\n" + 
    "    a    =          1\n" + 
    "   a =    1\n" + 
    "a = 5\n"

    p := New(lexer.New(input))

    if p.l.Indents != 4 {
        t.Errorf("expect %v, got %v", 4, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 4 {
        t.Errorf("expect %v, got %v", 4, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 8 {
        t.Errorf("expect %v, got %v", 8, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 4 {
        t.Errorf("expect %v, got %v", 4, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 3 {
        t.Errorf("expect %v, got %v", 4, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

}

func TestIndent2Parsing(t *testing.T) {
    input := "" +
    "a = 1\n" + 
    "if 2 > 1:\n" + 
    "    if 3 > 2:\n" + 
    "        a = 4\n" + 
    "        b = 4\n" + 
    "a = 5\n"

    p := New(lexer.New(input))

    if p.l.Indents != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }
}

