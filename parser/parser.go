package parser

import (
	"fmt"
	"os"
	"spl/lexer"
)

type Parser struct {
	tokens           []lexer.Token
	pos              int
	bracketCounter   int
	isInsideFunction bool
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0, bracketCounter: 0, isInsideFunction: false}
}

func (p *Parser) Parse() *Node {
	program := &Node{Type: "Program", Children: []*Node{}}
	for !p.match(lexer.EOF) {
		statement := p.parseStatement()
		program.Children = append(program.Children, statement)
		if p.match(lexer.RBRACKET) && p.bracketCounter > 0 {
			p.bracketCounter--
			p.moveRight()
			continue
		}
		if p.match(lexer.SEMICOLON) {
			p.moveRight()
		} else {
			panic("Expected terminator \";\"")
		}
	}
	return program
}

func (p *Parser) parseStatement() *Node {
	if p.match(lexer.IDENT) {
		p.moveRight()

		if p.match(lexer.ASSIGN) {
			p.moveLeft()
			return p.parseAssignmentStatement()
		}
	}

	if p.match(lexer.FUNC) && !p.isInsideFunction {
		return p.parseFunction()
	}

	if p.match(lexer.PRINT) {
		return p.ParsePrintStatement()
	}

	if p.match(lexer.LBRACKET) {
		p.bracketCounter++
		return p.parseScope()
	}

	if p.match(lexer.RETURN) && p.isInsideFunction {
		return p.parseReturnStatement()
	}

	panic("ParseStatement error")
}

func (p *Parser) ParsePrintStatement() *Node {
	p.moveRight()
	if p.match(lexer.SEMICOLON) {
		return &Node{Type: "PrintStatement", Children: []*Node{}}
	}
	expression := p.parseExpression()
	return &Node{Type: "PrintStatement", Children: []*Node{expression}}
}

func (p *Parser) parseFunction() *Node {
	identifier := p.moveRight()
	if p.match(lexer.LPAREN) {
		p.moveRight()
		parameters := p.parseParameterList()
		if p.match(lexer.COLON) {
			for _, value := range parameters.Children {
				if value.Type == lexer.INT.String() || value.Type == lexer.FLOAT.String() {
					panic("parameters in functions definitions must be identifiers, not numbers")
				}
			}
			p.moveRight()
			if p.match(lexer.LBRACKET) {
				p.isInsideFunction = true
				p.bracketCounter++
				function := &Node{Type: "FunctionDefinition", Children: []*Node{identifier, parameters, p.parseScope()}}
				p.isInsideFunction = false
				return function
			}
			defenition := p.parseFunctionDefinition()
			function := &Node{Type: "FunctionDefinition", Children: []*Node{identifier, parameters, defenition}}
			return function
		}
		return &Node{Type: "FunctionCall", Children: []*Node{identifier, parameters}}
	}
	panic("Expected '(' after function identifier")
}

func (p *Parser) parseFunctionDefinition() *Node {
	defenition := p.parseExpression()
	return &Node{Type: "ReturnStatement", Value: "", Children: []*Node{defenition}}

}

func (p *Parser) parseParameterList() *Node {
	parameters := &Node{Type: "ParameterList", Children: []*Node{}}
	for !p.match(lexer.RPAREN) {
		expression := p.parseExpression()
		parameters.Children = append(parameters.Children, expression)
		if p.match(lexer.COMMA) {
			p.moveRight()
			continue
		}
		if !p.match(lexer.RPAREN) {
			panic("expected ',' or '(' in parameter list")
		}
	}
	p.moveRight()
	return parameters
}

func (p *Parser) parseReturnStatement() *Node {
	p.moveRight()
	return &Node{Type: "ReturnStatement", Children: []*Node{p.parseExpression()}}
}

func (p *Parser) parseAssignmentStatement() *Node {
	identifier := p.moveRight()
	p.moveRight()
	expression := p.parseExpression()
	return &Node{Type: "AssignmentStatement", Value: identifier.Value, Children: []*Node{expression}}
}

func (p *Parser) parseExpression() *Node {
	node := p.parseTerm()

	for p.match(lexer.ADD) || p.match(lexer.SUB) {
		op := p.moveRight()
		child := p.parseTerm()
		node = &Node{Type: "BinaryOp", Value: op.Value, Children: []*Node{node, child}}
	}

	return node
}

func (p *Parser) parseTerm() *Node {
	node := p.parseFactor()

	for p.match(lexer.MUL) || p.match(lexer.DIV) {
		op := p.moveRight()
		child := p.parseFactor()
		node = &Node{Type: "BinaryOp", Value: op.Value, Children: []*Node{node, child}}
	}

	return node
}

func (p *Parser) parseFactor() *Node {
	if p.match(lexer.INT) || p.match(lexer.FLOAT) || p.match(lexer.IDENT) {
		node := p.moveRight()
		return node
	}

	if p.match(lexer.FUNC) && !p.isInsideFunction {
		return p.parseFunction()
	}

	if p.match(lexer.LPAREN) {
		p.moveRight()
		node := p.parseExpression()
		if !p.match(lexer.RPAREN) {
			panic("expected closing parenthesis")
		}
		p.moveRight()
		return node
	}

	panic("expected factor")
}

func (p *Parser) parseScope() *Node {
	block := &Node{Type: "Block", Children: []*Node{}}
	p.moveRight()
	for !p.match(lexer.EOF) {
		if p.match(lexer.RBRACKET) {
			return block
		}

		statement := p.parseStatement()
		block.Children = append(block.Children, statement)
		if p.match(lexer.RBRACKET) && p.bracketCounter > 0 {
			p.bracketCounter--
			p.moveRight()
			continue
		}
		if p.match(lexer.SEMICOLON) {
			p.moveRight()
		} else {
			fmt.Println("error: expected terminator ';'")
			os.Exit(3)
		}
	}
	p.moveRight()
	return block
}

func (p *Parser) match(expected lexer.TokenType) bool {
	return p.tokens[p.pos].Type == expected
}

func (p *Parser) moveRight() *Node {
	token := p.tokens[p.pos]
	p.pos++
	return &Node{Type: token.Type.String(), Value: token.Value}
}

func (p *Parser) moveLeft() *Node {
	token := p.tokens[p.pos]
	p.pos--
	return &Node{Type: token.Type.String(), Value: token.Value}
}

func (p *Parser) Clear() {
	p.pos = 0
	p.tokens = nil
}
