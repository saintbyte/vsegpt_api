package vsegpt

// Цена за токен у модели
type ModelPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

// Модель
type ModelItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ModelPricing  `json:"pricing"`
	ContextLength string `json:"context_length"`
}

// Формат ответа сервера со списком моделей
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelItem `json:"data"`
}

// Сообщение в запросе
type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Статистика запроса
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

// Сообщение в ответе.
type MessageResponse struct {
	Role           string     `json:"role"`
	Content        string     `json:"content"`
	DataForContext []struct{} `json:"data_for_context"`
}

// Варианты.
type ChoicesResponse struct {
	Message      MessageRequest `json:"message"`
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
}

// Исходящий запрос к модели чата.
type ChatCompletionRequest struct {
	Model             string           `json:"model"`
	Messages          []MessageRequest `json:"messages"`
	Stream            bool             `json:"stream"`
	RepetitionPenalty float32          `json:"repetition_penalty,omitempty"`
	Temperature       float32          `json:"temperature,omitempty"`
	N                 int              `json:"n,omitempty"`
	TopP              float32          `json:"top_p,omitempty"`
	MaxTokens         int              `json:"max_tokens,omitempty"`
	UpdateInterval    int              `json:"update_interval,omitempty"`
}

// Ответ от модели на запрос.
type ChatCompletionResponse struct {
	Choices []ChoicesResponse `json:"choices"`
	Created int               `json:"created"`
	Model   string            `json:"model"`
	Usage   Usage             `json:"usage"`
	Object  string            `json:"object"`
}

// Запрос на векторизацию текста
type EmbeddingsRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"` // default float
}

// Ответ с векторизацией текста
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
