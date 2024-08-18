package vsegpt

type ModelPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}
type ModelItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ModelPricing  `json:"pricing"`
	ContextLength string `json:"context_length"`
}

type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelItem `json:"data"`
}

type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type MessageResponse struct {
	Role           string     `json:"role"`
	Content        string     `json:"content"`
	DataForContext []struct{} `json:"data_for_context"`
}

type ChoicesResponse struct {
	Message      MessageRequest `json:"message"`
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
}

type ChatCompletionRequest struct {
	Model             string           `json:"model"`
	Messages          []MessageRequest `json:"messages"`
	Stream            bool             `json:"stream"`
	RepetitionPenalty int              `json:"repetition_penalty,omitempty"`
	Temperature       float32          `json:"temperature,omitempty"`
	TopP              float32          `json:"top_p,omitempty"`
	MaxTokens         int              `json:"max_tokens,omitempty"`
	UpdateInterval    int              `json:"update_interval,omitempty"`
}

type ChatCompletionResponse struct {
	Choices []ChoicesResponse `json:"choices"`
	Created int               `json:"created"`
	Model   string            `json:"model"`
	Usage   Usage             `json:"usage"`
	Object  string            `json:"object"`
}

type EmbeddingsRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"` // default float
}

type EmbeddingsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage Usage  `json:"usage"`
}
