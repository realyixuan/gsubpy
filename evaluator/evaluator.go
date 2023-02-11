package evaluator

import (
    "gsubpy/ast"
    "gsubpy/object"
)

func NewEnvironment() *Environment {
    s := make(map[string]object.Object)
    return &Environment{store: s}
}

type Environment struct {
    store map[string]object.Object
}

func (e *Environment) Get(name string) (object.Object, bool) {
    obj, ok := e.store[name]
    return obj, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
    e.store[name] = val
    return val
}

func Eval(node ast.Node, env *Environment) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node, env)
    // case *ast.ExpressionStatement:
    //     return Eval(node.Expression)
    // case *ast.IntegerLiteral:
    //     return &object.Integer{Value: node.Value}
    case *ast.AssignmentStatement:
        val := Eval(node.Value, env)
        env.Set(node.Name.Value, val)
    // case *ast.Identifier:
    //     return evalIdentifier(node, env)

    case *ast.InfixExpression:
        left := Eval(node.Left, env)
        right := Eval(node.Right, env)
        return evalInfixExpression(node.Operator, left, right)

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}


    }

    return nil
}

func evalInfixExpression(
    operator string,
    left, right object.Object,
) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        leftVal := left.(*object.Integer).Value
        rightVal := right.(*object.Integer).Value
        if operator == "+" {
            return &object.Integer{Value: leftVal + rightVal}
        } else if operator == "-" {
            return &object.Integer{Value: leftVal - rightVal}
        }
    }
    return nil
}

func evalIdentifier(
    node *ast.Identifier,
    env *Environment,
) object.Object {
    val, _ := env.Get(node.Value)
    return val
}

func evalProgram(program *ast.Program, env *Environment) object.Object {
    var result object.Object
    for _, statement := range program.Statements {
        result = Eval(statement, env)
    }

    return result
}

// func evalStatements(stmts []ast.Statement) object.Object {
//     var result object.Object
// 
//     for _, statement := range stmts {
//         result = Eval(statement)
//     }
// 
//     return result
// }

