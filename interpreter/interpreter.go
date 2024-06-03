package interpreter

import (
	"fmt"
	"os"
	"spl/lexer"
	"spl/parser"
	"strconv"
)

type Interpreter struct {
	currentScope *Scope
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		currentScope: NewScope(nil),
	}
}

func (i *Interpreter) EnterScope() {
	i.currentScope = NewScope(i.currentScope)
}

func (i *Interpreter) ExitScope() {
	if i.currentScope != nil {
		i.currentScope = i.currentScope.parent
	} else {
		fmt.Println("error: attempt to exit global scopa")
		os.Exit(1)
	}

}

func (i *Interpreter) AssignVariable(name, value string) {
	i.currentScope.vars[name] = value
}

func (i *Interpreter) FindVariable(name string) (string, error) {
	scope := i.currentScope
	for scope != nil {
		if value, ok := scope.vars[name]; ok {
			return value, nil
		}
		scope = scope.parent
	}
	return "", fmt.Errorf("variable %s not defined", name)
}

func (i *Interpreter) AssignFunction(name string, node *parser.Node) {
	i.currentScope.functions[name] = node
}

func (i *Interpreter) FindFunction(name string) (*parser.Node, error) {
	scope := i.currentScope
	for scope != nil {
		if function, ok := scope.functions[name]; ok {
			return function, nil
		}
		scope = scope.parent
	}
	return nil, fmt.Errorf("function %s not defined", name)
}

func (i *Interpreter) Execute(node *parser.Node) error {
	for len(node.Children) != 0 {
		switch node.Children[0].Type {
		case "AssignmentStatement":
			varName := node.Children[0].Value
			value, err := i.Evaluate(node.Children[0].Children[0])
			if err != nil {
				return err
			}
			i.AssignVariable(varName, value)
		case "PrintStatement":
			switch len(node.Children[0].Children) == 0 {
			case true:
				i.printAllVars()
			case false:
				value, err := i.Evaluate(node.Children[0].Children[0])
				if err != nil {
					return err
				}
				fmt.Println(value)
			}
		case "FunctionDefinition":
			// i.functions[node.Children[0].Children[0].Value] = node.Children[0]
			i.AssignFunction(node.Children[0].Children[0].Value, node.Children[0])
		case "Block":
			i.EnterScope()
			i.Execute(node.Children[0])
			i.ExitScope()
		default:
			_, err := i.Evaluate(node.Children[0])
			if err != nil {
				return err
			}
		}

		node.Children = node.Children[1:]
	}
	return nil
}

func (i *Interpreter) Evaluate(node *parser.Node) (string, error) {
	switch node.Type {
	case lexer.INT.String(), lexer.FLOAT.String():
		return node.Value, nil
	case lexer.IDENT.String():
		return i.FindVariable(node.Value)
	case "BinaryOp":
		leftValue, err := i.Evaluate(node.Children[0])
		if err != nil {
			return "", fmt.Errorf("left value error")
		}

		rightValue, err := i.Evaluate(node.Children[1])
		if err != nil {
			return "", fmt.Errorf("right value error")
		}

		return applyOperator(leftValue, rightValue, node.Value)

	case "FunctionCall":
		return i.evaluateFunctionCall(node)
	default:
		return "", fmt.Errorf("unexpected token type: %v", node.Type)
	}
}

func applyOperator(left, right, operator string) (string, error) {
	leftInt, errLeftInt := strconv.Atoi(left)
	rightInt, errRightInt := strconv.Atoi(right)

	if errLeftInt == nil && errRightInt == nil {
		var result int
		switch operator {
		case "+":
			result = leftInt + rightInt
		case "-":
			result = leftInt - rightInt
		case "*":
			result = leftInt * rightInt
		case "/":
			if rightInt == 0 {
				return "", fmt.Errorf("division by zero")
			}
			result = leftInt / rightInt
		default:
			return "", fmt.Errorf("unknow operator %s", operator)
		}
		return strconv.FormatInt(int64(result), 10), nil
	}

	leftFloat, errLeftFloat := strconv.ParseFloat(left, 64)
	rightFloat, errRightFloat := strconv.ParseFloat(right, 64)

	if errLeftFloat != nil || errRightFloat != nil {
		return "", fmt.Errorf("invalid operands: %s, %s", left, right)
	}

	var result float64
	switch operator {
	case "+":
		result = leftFloat + rightFloat
	case "-":
		result = leftFloat - rightFloat
	case "*":
		result = leftFloat * rightFloat
	case "/":
		if rightFloat == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = leftFloat / rightFloat
	default:
		return "", fmt.Errorf("unknow operator %s", operator)
	}
	if result == float64(int(result)) {
		return strconv.FormatFloat(result, 'f', 1, 64), nil
	}
	return strconv.FormatFloat(result, 'f', -1, 64), nil

}

func (i *Interpreter) evaluateFunctionCall(node *parser.Node) (string, error) {
	functionName := node.Children[0].Value
	function, err := i.FindFunction(functionName)
	if err != nil {
		return "", err
	}

	parameters := function.Children[1].Children
	arguments := node.Children[1].Children

	if len(parameters) != len(arguments) {
		return "", fmt.Errorf("incorrect number of arguments for function %s", functionName)
	}

	savedVars := make(map[string]string)
	for key, value := range i.currentScope.vars {
		savedVars[key] = value
	}

	for index, param := range parameters {
		argValue, err := i.Evaluate(arguments[index])
		if err != nil {
			return "", err
		}

		i.currentScope.vars[param.Value] = argValue
	}

	result, err := i.Evaluate(function.Children[2].Children[0])
	if err != nil {
		return "", err
	}

	i.currentScope.vars = savedVars

	return result, nil
}

func (i *Interpreter) printAllVars() {
	for key, value := range i.currentScope.vars {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func (i *Interpreter) Clear() {
	i.currentScope.vars = nil
	i.currentScope.functions = nil
}
