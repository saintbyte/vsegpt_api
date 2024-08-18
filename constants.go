package vsegpt

const (
	VseGptApiKeyEnv          = "VSEGPT_API_KEY"
	VseGptApiHost            = "api.vsegpt.ru" //api.vsegpt.ru:6070
	VseGptModelsPath         = "/v1/models"
	VseGptChatCompletionPath = "/v1/chat/completions"
	VseGptEmbeddingsPath     = "/v1/embeddings"
	VseGptEmbeddingModel     = "emb-openai/text-embedding-3-small"
	VseGptModel              = "openai/gpt-3.5-turbo"
	VseGptMaxTokens          = 16384
	MaxEmbeddingSize         = 8192
)
