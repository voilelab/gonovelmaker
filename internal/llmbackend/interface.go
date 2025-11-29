package llmbackend

import "context"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

type LLMBackend interface {
	ChatCompletion(messages []Message, ctx context.Context) (string, error)
}
