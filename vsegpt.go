package vsegpt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

type VseGpt struct {
	ApiKey           string
	Model            string
	MaxTokens        int
	MaxEmbeddingSize int
	EmbeddingModel   string
}

func NewVseGpt() *VseGpt {
	return &VseGpt{
		ApiKey:           "",
		Model:            VseGptModel,
		MaxTokens:        VseGptMaxTokens,
		MaxEmbeddingSize: 8192,
		EmbeddingModel:   VseGptEmbeddingModel,
	}
}
func (v *VseGpt) GetRequestUrl(path string) string {
	return "https://" + VseGptApiHost + path
}

func (v *VseGpt) GetRequest(url string) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+v.GetCurrentToken())
	return request, nil
}

func (v *VseGpt) PostRequest(url string, body io.Reader) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+v.GetCurrentToken())
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (v *VseGpt) GetCurrentToken() string {
	value, exists := os.LookupEnv(VseGptApiKeyEnv)
	if exists {
		return value
	}
	if v.ApiKey != "" {
		return v.ApiKey
	}
	return ""
}

func (v *VseGpt) GetModels() ([]ModelItem, error) {
	url := v.GetRequestUrl(VseGptModelsPath)
	request, err := v.GetRequest(url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Response status: " + string(response.Status))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result ModelsResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	return result.Data, nil
}

func (v *VseGpt) Embeddings(input string) ([]float64, error) {
	url := v.GetRequestUrl(VseGptEmbeddingsPath)
	jData, errJsonRequestEncode := json.Marshal(&EmbeddingsRequest{
		Model:          v.EmbeddingModel,
		Input:          input,
		EncodingFormat: "float",
	})
	if errJsonRequestEncode != nil {
		return nil, errJsonRequestEncode
	}
	request, err := v.PostRequest(url, bytes.NewReader(jData))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		return nil, errors.New("Response status: " + string(response.Status) + " " + string(body))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result EmbeddingsResponse
	err = json.Unmarshal(body, &result)
	return result.Data[0].Embedding, nil
}

func (v *VseGpt) ChatCompletion(messages []MessageRequest) (string, error) {
	url := v.GetRequestUrl(VseGptChatCompletionPath)
	jData, errJsonRequestEncode := json.Marshal(&ChatCompletionRequest{
		Model:    v.Model,
		Messages: messages,
		Stream:   false,
	})
	if errJsonRequestEncode != nil {
		return "", errJsonRequestEncode
	}
	request, err := v.PostRequest(url, bytes.NewReader(jData))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		return "", errors.New("Response status: " + string(response.Status) + " " + string(body))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result ChatCompletionResponse
	err = json.Unmarshal(body, &result)
	return result.Choices[0].Message.Content, nil
}

func (v *VseGpt) Ask(question string) (string, error) {
	return v.ChatCompletion([]MessageRequest{
		MessageRequest{
			Role:    "user",
			Content: question,
		},
	})
}
