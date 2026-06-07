package domain

type Role string

const (
	RoleUser Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role Role
	Blocks []Block
}

type Block interface {
	isBlock()
}

type TextBlock struct {
	Text string
}
func (TextBlock) isBlock() {}

