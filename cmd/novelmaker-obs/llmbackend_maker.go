package main

import "github.com/voilelab/gonovelmaker/internal/llmbackend"

type LLMBackendMaker func(apiKey, model, baseURL string) llmbackend.LLMBackend

func openAIBackendMaker(apiKey, model, baseURL string) llmbackend.LLMBackend {
	return llmbackend.NewOpenAIBackend(
		apiKey,
		model,
		baseURL,
	)
}

func dummyBackendMaker(_, _, _ string) llmbackend.LLMBackend {
	return llmbackend.NewDummyBackend()
}
