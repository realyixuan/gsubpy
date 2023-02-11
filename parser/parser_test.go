package parser

import (
    "testing"
    "gsubpy/ast"
    "gsubpy/lexer"
)

func TestAssignmentStatements(t *testing.T) {
    input := `sum = 11`
    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    if program == nil {
        t.Fatalf("parseProgram() returned nil")
    }

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain 1 statements. got=%d",
            len(program.Statements))
    }
}

func testAssignmentStatement(t *testing.T, s ast.Statement, name string) bool {
    // if s.TokenLiteral() != ""
    AssignmentStmt, ok := s.(*ast.AssignmentStatement)
    if !ok {
        t.Errorf("s not *ast.AssignmentStatement. got=%T", s)
        return false
    }

    if AssignmentStmt.Name.Value != name {
        t.Errorf("AssignmentStmt.Name.Value not '%s'. got=%T", name, AssignmentStmt.Name.Value)
        return false
    }

    if AssignmentStmt.Name.TokenLiteral() != name {
        t.Errorf("s.Name not '%s'. got=%T", name, AssignmentStmt.Name)
        return false
    }

    return true
}

