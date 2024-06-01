package main

import (
	"fmt"
	"spl/interpreter"
	"spl/lexer"
	"spl/parser"
	"strings"
)

func main() {
	input := `
	foo(x, y): ((x*y+2)*(25-x/y))/(3-(x+2*y));
	myfoo2(z): z*z+4;
	myvar=15;
	bg=25.0;
	ccc=myfoo2(bg+myvar)*15+foo(bg*25,(6*myfoo2(myvar-10)));
	print ccc;
	bg=ccc*myvar;
	print;
	`
	var tokens []lexer.Token
	lex := lexer.NewLexer(input)
	for {
		tok := lex.NextToken()
		tokens = append(tokens, tok)
		// fmt.Printf("%+v\n", tok)
		if tok.Type == lexer.EOF {
			break
		}
	}
	parser := parser.NewParser(tokens)
	ast := parser.Parse()

	// // Вывод AST в виде дерева
	// printAST(ast, 0)
	inter := interpreter.NewInterpreter()
	err := inter.Execute(ast)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Функция для вывода AST в виде дерева
func printAST(node *parser.Node, indent int) {
	fmt.Println(strings.Repeat("*", indent*4), node.Type, node.Value)
	for _, child := range node.Children {
		printAST(child, indent+1)
	}
}
