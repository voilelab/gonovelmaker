package llmbackend

type LLMBackendMaker func(apiKey, baseURL, model string) LLMBackend

func MakeOpenAI(apiKey, baseURL, model string) LLMBackend {
	return NewOpenAIBackend(
		apiKey,
		baseURL,
		model,
	)
}

func MakeDummy(_, _, _ string) LLMBackend {
	return NewDummyBackend()
}
