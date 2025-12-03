package llmbackend

import "context"

var _ LLMBackend = (*DummyBackend)(nil)

type DummyBackend struct{}

func NewDummyBackend() *DummyBackend {
	return &DummyBackend{}
}

func (d *DummyBackend) ChatCompletion(messages []Message, ctx context.Context) (string, UsageInfo, error) {
	return "This is a dummy response from DummyBackend.", UsageInfo{}, nil
}

func (d *DummyBackend) GenerateImage(prompt string, ctx context.Context) (string, error) {
	return "https://example.com/dummy-image.png", nil
}
