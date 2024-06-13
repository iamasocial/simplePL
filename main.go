package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"spl/interpreter"
	"spl/lexer"
	"spl/parser"
	"strings"
)

type Cleaner interface {
	Clear()
}

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var input string
	for scanner.Scan() {
		input += scanner.Text()
	}
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
	_, err = inter.Execute(ast)
	if err != nil {
		fmt.Println(err)
		return
	}
	clearMemory(parser, inter)
	runtime.GC()
}

// Функция для вывода AST в виде дерева
func printAST(node *parser.Node, indent int) {
	fmt.Println(strings.Repeat("*", indent*4), node.Type, node.Value)
	for _, child := range node.Children {
		printAST(child, indent+1)
	}
}

func clearMemory(Cleaner ...Cleaner) {
	for _, object := range Cleaner {
		object.Clear()
	}
}
