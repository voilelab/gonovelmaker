package main

import "github.com/voilelab/gonovelmaker/internal/llmbackend"

type LLMBackendMaker func(apiKey, baseURL, model string) llmbackend.LLMBackend

func openAIBackendMaker(apiKey, baseURL, model string) llmbackend.LLMBackend {
	return llmbackend.NewOpenAIBackend(
		apiKey,
		baseURL,
		model,
	)
}

func dummyBackendMaker(_, _, _ string) llmbackend.LLMBackend {
	return llmbackend.NewDummyBackend()
}
