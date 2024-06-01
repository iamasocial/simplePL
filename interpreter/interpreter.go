package interpreter

import (
	"fmt"
	"spl/lexer"
	"spl/parser"
	"strconv"
)

type Interpreter struct {
	vars      map[string]string
	functions map[string]*parser.Node
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		vars:      make(map[string]string),
		functions: make(map[string]*parser.Node),
	}
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
			i.vars[varName] = value
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
			i.functions[node.Children[0].Children[0].Value] = node.Children[0]
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
	case lexer.INT.String():
		return node.Value, nil
	case lexer.FLOAT.String():
		return node.Value, nil
	case lexer.IDENT.String():
		value, ok := i.vars[node.Value]
		if !ok {
			return "", fmt.Errorf("variable %s not defined", node.Value)
		}
		return value, nil
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
	function, ok := i.functions[functionName]
	if !ok {
		return "", fmt.Errorf("function %s not defined", functionName)
	}

	parameters := function.Children[1].Children
	arguments := node.Children[1].Children

	if len(parameters) != len(arguments) {
		return "", fmt.Errorf("incorrect number of arguments for function %s", functionName)
	}

	savedVars := make(map[string]string)
	for key, value := range i.vars {
		savedVars[key] = value
	}

	for idx, param := range parameters {
		argValue, err := i.Evaluate(arguments[idx])
		if err != nil {
			return "", err
		}

		i.vars[param.Value] = argValue
	}

	result, err := i.Evaluate(function.Children[2].Children[0])
	if err != nil {
		return "", err
	}

	i.vars = savedVars

	return result, nil
}

func (i *Interpreter) printAllVars() {
	for key, value := range i.vars {
		fmt.Printf("%s: %s\n", key, value)
	}
}
