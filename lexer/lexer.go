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
	for l.currentChar == ' ' || l.currentChar == '\n' || l.currentChar == '\r' || l.currentChar == '\t' {
		l.readChar()
	}

	// Проверка на конец файла
	if l.currentChar == 0 {
		tok = NewToken(EOF, "")
	} else if l.currentChar == '(' {
		tok = NewToken(LPAREN, string(l.currentChar))
	} else if l.currentChar == ')' {
		tok = NewToken(RPAREN, string(l.currentChar))
	} else if l.currentChar == ',' {
		tok = NewToken(COMMA, string(l.currentChar))
	} else if l.currentChar == '+' {
		tok = NewToken(ADD, string(l.currentChar))
	} else if l.currentChar == '-' {
		tok = NewToken(SUB, string(l.currentChar))
	} else if l.currentChar == '*' {
		tok = NewToken(MUL, string(l.currentChar))
	} else if l.currentChar == '/' {
		tok = NewToken(DIV, string(l.currentChar))
	} else if l.currentChar == '=' {
		tok = NewToken(ASSIGN, string(l.currentChar))
	} else if l.currentChar == ':' {
		tok = NewToken(COLON, string(l.currentChar))
	} else if l.currentChar == ';' {
		tok = NewToken(SEMICOLON, string(l.currentChar))
	} else if unicode.IsDigit(rune(l.currentChar)) {
		// tok.Type = NUMBER
		// var num strings.Builder
		// for unicode.IsDigit(rune(l.currentChar)) {
		// 	num.WriteByte(l.currentChar)
		// 	l.readChar()
		// }
		// l.Position--
		// tok.Value = num.String()
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
		l.Position--
	} else if unicode.IsLetter(rune(l.currentChar)) {
		var ident strings.Builder
		for unicode.IsLetter(rune(l.currentChar)) || unicode.IsDigit(rune(l.currentChar)) {
			ident.WriteByte(l.currentChar)
			l.readChar()
		}
		if l.currentChar == '(' {
			tok = NewToken(FUNC, ident.String())
			return tok
		}
		l.Position--
		switch ident.String() {
		case "print":
			tok = NewToken(PRINT, ident.String())
		default:
			tok = NewToken(IDENT, ident.String())
		}
	} else {
		tok = NewToken(ILLEGAL, string(l.currentChar))
		panic(fmt.Sprintf("Found illegal token: %s", tok.Value))
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
