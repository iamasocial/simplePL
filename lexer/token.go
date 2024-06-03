package lexer

type TokenType int

const (
	IDENT TokenType = iota
	NUMBER
	INT
	FLOAT
	FUNC
	PRINT
	ASSIGN
	ADD
	SUB
	MUL
	DIV
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	COMMA
	COLON
	SEMICOLON
	EOF
	ILLEGAL
)

var tokens = []string{
	IDENT:     "IDENT",
	NUMBER:    "NUMBER",
	INT:       "INT",
	FLOAT:     "FLOAT",
	FUNC:      "FUNC",
	PRINT:     "PRINT",
	ASSIGN:    "ASSIGN",
	ADD:       "ADD",
	SUB:       "SUB",
	MUL:       "MUL",
	DIV:       "DIV",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACKET:  "LBRACKET",
	RBRACKET:  "RBRACKET",
	COMMA:     "COMMA",
	COLON:     "COLON",
	SEMICOLON: "SEMICOLON",
	EOF:       "EOF",
	ILLEGAL:   "ILLGEGAL",
}

func (t TokenType) String() string {
	return tokens[t]
}

// Структура для токена
type Token struct {
	Type  TokenType
	Value string
}

// Функция для создания токена
func NewToken(tokenType TokenType, value string) Token {
	return Token{Type: tokenType, Value: value}
}
