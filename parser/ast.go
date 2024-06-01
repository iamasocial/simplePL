package parser

// Структура AST узла
type Node struct {
	Type     string
	Value    string
	Children []*Node
}

// Функция для создания узла AST
func NewNode(nodeType, value string) *Node {
	return &Node{Type: nodeType, Value: value}
}
