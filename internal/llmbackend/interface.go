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

type UsageInfo struct {
	InputTokens  int64 // Tokens used in the prompt
	OutputTokens int64 // Tokens generated in the completion
	TotalTokens  int64 // Total tokens used
}

type LLMBackend interface {
	ChatCompletion(messages []Message, ctx context.Context) (content string, usage UsageInfo, err error)
	GenerateImage(prompt string, ctx context.Context) (imageURL string, err error)
	ListModels(ctx context.Context) ([]string, error)
}
