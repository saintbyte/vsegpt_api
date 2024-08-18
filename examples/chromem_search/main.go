package main

import (
	"context"
	"github.com/philippgille/chromem-go"
	chromem_indexer "github.com/saintbyte/golang_gigachat_api/pkg/chromem"
	"github.com/saintbyte/golang_gigachat_api/pkg/vsegpt"
	"log"
	"strings"
	"time"
)

const (
	MinSimilarity float32 = 0.4
)

func main() {
	ctx := context.Background()
	log.Println("Setting up chromem-go...")
	db, err := chromem.NewPersistentDB("./db", false)
	if err != nil {
		panic(err)
	}
	collection := db.GetCollection("civil_codex", chromem_indexer.NewEmbeddingFuncVseGpt())
	if err != nil {
		panic(err)
	}
	start := time.Now()
	log.Println("Querying chromem-go...")
	// "nomic-embed-text" specific prefix (not required with OpenAI's or other models)
	question := "Расскажи про договоры?"
	query := "search_query: " + question
	docRes, err := collection.Query(ctx, query, 2, nil, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Search (incl query embedding) took", time.Since(start))
	// Here you could filter out any documents whose similarity is below a certain threshold.
	// if docRes[...].Similarity < 0.5 { ...

	// Print the retrieved documents and their similarity to the question.
	for i, res := range docRes {
		// Cut off the prefix we added before adding the document (see comment above).
		// This is specific to the "nomic-embed-text" model.
		content := strings.TrimPrefix(res.Content, "search_document: ")
		log.Printf("Document %d (similarity: %f): \"%s\"\n", i+1, res.Similarity, content)
	}
	goodQuestion := false
	if (len(docRes) > 0) && (docRes[0].Similarity >= MinSimilarity) {
		goodQuestion = true
	}
	// Now we can ask the LLM again, augmenting the question with the knowledge we retrieved.
	// In this example we just use both retrieved documents as context.
	//contexts := []string{docRes[0].Content, docRes[1].Content}
	log.Println("Asking LLM with augmented question...")
	if goodQuestion {
		gpt35turbo := vsegpt.NewVseGpt()
		messages := []vsegpt.MessageRequest{
			vsegpt.MessageRequest{
				Role:    "system",
				Content: docRes[0].Content,
			},
			vsegpt.MessageRequest{
				Role:    "system",
				Content: docRes[1].Content,
			},
			vsegpt.MessageRequest{

				Role:    "user",
				Content: question,
			},
		}
		aswer, _ := gpt35turbo.ChatCompletion(vsegpt.VseGptModel, messages)
		log.Println(aswer)
	} else {
		log.Println("Что-то я не похоже на вопрос по гражданскому праву.")
	}
	///log.Printf("Reply after augmenting the question with knowledge: \"" + reply + "\"\n")
}
