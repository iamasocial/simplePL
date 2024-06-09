package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Лексер
type Lexer struct {
	input       string
	Position    int
	currentChar byte
	keywords    map[string]TokenType
}

// Инициализация лексера
func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input, keywords: map[string]TokenType{
		"print": PRINT,
		// "func":  FUNC,
	}}
	lexer.readChar()
	return lexer
}

// Чтение очередного символа
func (l *Lexer) readChar() {
	if l.Position < len(l.input) {
		l.currentChar = l.input[l.Position]
	} else {
		l.currentChar = 0 // ASCII код NUL
	}
	l.Position++
}

// Получение следующего токена
func (l *Lexer) NextToken() Token {
	var tok Token

	// Пропускаем пробелы и символы перевода строки
	l.skipSpace()
	switch l.currentChar {
	case 0:
		tok = NewToken(EOF, "")
	case '(':
		tok = NewToken(LPAREN, "(")
	case ')':
		tok = NewToken(RPAREN, ")")
	case '{':
		tok = NewToken(LBRACKET, "{")
	case '}':
		tok = NewToken(RBRACKET, "}")
	case ',':
		tok = NewToken(COMMA, ",")
	case '+':
		tok = NewToken(ADD, "+")
	case '-':
		tok = NewToken(SUB, "-")
	case '*':
		tok = NewToken(MUL, "*")
	case '/':
		tok = NewToken(DIV, "/")
	case '=':
		tok = NewToken(ASSIGN, "=")
	case ':':
		tok = NewToken(COLON, ":")
	case ';':
		tok = NewToken(SEMICOLON, ";")
	default:
		if unicode.IsDigit(rune(l.currentChar)) {
			num, err := l.lexNumber()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if strings.Contains(num, ".") {
				tok = NewToken(FLOAT, num)
			} else {
				tok = NewToken(INT, num)
			}
			return tok
		}

		if unicode.IsLetter(rune(l.currentChar)) {
			var ident strings.Builder
			for unicode.IsLetter(rune(l.currentChar)) || unicode.IsDigit(rune(l.currentChar)) {
				ident.WriteByte(l.currentChar)
				l.readChar()
			}

			if l.currentChar == '(' {
				tok = NewToken(FUNC, ident.String())
				return tok
			}

			switch ident.String() {
			case "print":
				tok = NewToken(PRINT, "print")
			case "return":
				tok = NewToken(RETURN, "return")
			default:
				tok = NewToken(IDENT, ident.String())
			}
			return tok
		}
		tok = NewToken(ILLEGAL, string(l.currentChar))
		fmt.Printf("error: found illegal token %s\n", tok.Value)
		os.Exit(2)
		return tok
	}
	l.readChar()
	return tok
}

func (l *Lexer) lexNumber() (string, error) {
	var num strings.Builder
	dotCounter := 0
	for {
		if unicode.IsDigit(rune(l.currentChar)) {
			num.WriteByte(l.currentChar)
			l.readChar()
			continue
		}

		if l.currentChar == '.' {
			if dotCounter > 0 {
				return "", fmt.Errorf("float number can have only 1 dot")
			}

			num.WriteByte(l.currentChar)
			l.readChar()
			continue
		}

		return num.String(), nil
	}

}

func (l *Lexer) skipSpace() {
	for l.currentChar == ' ' || l.currentChar == '\n' || l.currentChar == '\r' || l.currentChar == '\t' {
		l.readChar()
	}
}
