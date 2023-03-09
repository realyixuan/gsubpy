package parser

import (
    "testing"

    "gsubpy/lexer"
)

func TestSpaceIndentParsing(t *testing.T) {
    input := "" +
    "a = 1\n" + 
    "if 2 > 1:\n" + 
    "    if 3 > 2:\n" + 
    "        a = 4\n" + 
    "        b = 4\n" + 
    "a = 5\n"

    p := New(lexer.New(input))

    if len(p.l.Indents) != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if len(p.l.Indents) != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if len(p.l.Indents) != 0 {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }
}

func TestTabIndentParsing(t *testing.T) {
    input := "" +
    "a = 1\n" + 
    "if 2 > 1:\n" + 
    "\tif 3 > 2:\n" + 
    "\t\ta = 4\n" + 
    "\t\tb = 4\n" + 
    "a = 5\n"

    p := New(lexer.New(input))

    if p.l.Indents != "" {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != "" {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != "" {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }
}

func TestTabSpaceIndentParsing(t *testing.T) {
    input := "" +
    "a = 1\n" + 
    "if 2 > 1:\n" + 
    "\t if 3 > 2:\n" + 
    "\t \t a = 4\n" + 
    "\t \t b = 4\n" + 
    "a = 5\n"

    p := New(lexer.New(input))

    if p.l.Indents != "" {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != "" {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }

    p.parsingStatement()
    if p.l.Indents != "" {
        t.Errorf("expect %v, got %v", 0, p.l.Indents)
    }
}

