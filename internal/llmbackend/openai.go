package llmbackend

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var _ LLMBackend = (*OpenAIBackend)(nil)

type OpenAIBackend struct {
	apiKey  string
	baseURL string
	model   string
}

func NewOpenAIBackend(apiKey, baseURL, model string) *OpenAIBackend {
	return &OpenAIBackend{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
	}
}

func (o *OpenAIBackend) ChatCompletion(messages []Message, ctx context.Context) (string, error) {
	msgs := []openai.ChatCompletionMessageParamUnion{}
	for _, m := range messages {
		switch m.Role {
		case RoleSystem:
			msgs = append(msgs, openai.SystemMessage(m.Content))
		case RoleUser:
			msgs = append(msgs, openai.UserMessage(m.Content))
		case RoleAssistant:
			msgs = append(msgs, openai.AssistantMessage(m.Content))
		default:
			return "", fmt.Errorf("unknown message role: %s", m.Role)
		}
	}

	opts := []option.RequestOption{option.WithAPIKey(o.apiKey)}
	if o.baseURL != "" {
		opts = append(opts, option.WithBaseURL(o.baseURL))
	}
	client := openai.NewClient(opts...)

	chatCompletion, err := client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Messages: msgs,
			Model:    o.model,
		})

	if err != nil {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
