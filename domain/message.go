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

// LLM がツールを呼ぼうとするときに帰ってくる
type ToolUseBlock struct {
	ID string
	Name string
	Input map[string]any // LLM が渡す引数
}

func (ToolUseBlock) isBlock() {}

// ツール実行結果を LLM に返す
type ToolResultBlock struct {
	ToolUseID string
	Content string
}

func (ToolResultBlock) isBlock() {}
