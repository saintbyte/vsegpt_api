package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/philippgille/chromem-go"
	chromem_indexer "github.com/saintbyte/golang_gigachat_api/pkg/chromem"
	"github.com/saintbyte/golang_gigachat_api/pkg/rag"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	log.Println("Setting up chromem-go...")
	db, err := chromem.NewPersistentDB("./db", false)
	if err != nil {
		panic(err)
	}
	collection, err := db.GetOrCreateCollection("civil_codex", nil, chromem_indexer.NewEmbeddingFuncVseGpt())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "501") {
			log.Println("Dawn 501 error")
			time.Sleep(10 * time.Second)
			collection, err = db.GetOrCreateCollection("civil_codex", nil, chromem_indexer.NewEmbeddingFuncVseGpt())
			if strings.Contains(fmt.Sprint(err), "501") {
				log.Println("Dawn 501 error !!!!!!")
				time.Sleep(30 * time.Second)
				collection, err = db.GetOrCreateCollection("civil_codex", nil, chromem_indexer.NewEmbeddingFuncVseGpt())

			}
		}
		panic(err)

	}
	if collection.Count() != 0 {
		log.Println("Collection already created")
		os.Exit(0)
	}
	f, err := os.Open("/home/sb/projects/civil_codex_gc/py_parser/pythonProject/result.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	d := json.NewDecoder(f)
	log.Println("Reading JSON lines...")
	//var docs []chromem.Document
	var records []rag.CivilCodexRecord
	err = d.Decode(&records)
	if err == io.EOF {
		//break // reached end of file
		log.Fatal("EOF")
	} else if err != nil {
		panic(err)
	}
	// Составляем массив
	counter := 0
	content := ""
	metadata := make(map[string]string)
	for i, record := range records {
		log.Println(record.Article, record.ArticleTitle)
		content = ""
		content = record.CodexName + " Часть " + string(record.Part)
		if record.SectionTitle != "" {
			content = content + " " + record.SectionTitle
		}
		if record.SubsectionTitle != "" {
			content = content + " " + record.SectionTitle
		}
		if record.ChapterTitle != "" {
			content = content + " Глава  " + record.Chapter + " " + record.ChapterTitle
		}
		content = content + " Статья  " + record.Article + " " + record.ArticleTitle
		content = content + " " + record.Text
		content = strings.TrimSpace(content)
		counter++
		metadata = map[string]string{
			"CodexName":       record.CodexName,
			"Part":            record.Part,
			"Section":         record.Section,
			"SectionTitle":    record.SectionTitle,
			"Subsection":      record.Subsection,
			"SubsectionTitle": record.SubsectionTitle,
			"Chapter":         record.Chapter,
			"ChapterTitle":    record.ChapterTitle,
			"Article":         record.Article,
			"ArticleTitle":    record.ArticleTitle,
			"Point":           record.Point,
		}
		err = collection.AddDocument(ctx,
			chromem.Document{
				ID:       string(i),
				Metadata: metadata,
				Content:  content,
			})
		if err != nil {
			log.Println("Add documents error: ")
			log.Println(err)
			log.Println("Some wait and try again")
			time.Sleep(20 * time.Second)
			err = collection.AddDocument(ctx,
				chromem.Document{
					ID:       string(i),
					Metadata: metadata,
					Content:  content,
				})
			if err != nil {
				log.Println(err)
				log.Println("Some wait and try again 2 !!!!")
				time.Sleep(20 * time.Second)
				err = collection.AddDocument(ctx,
					chromem.Document{
						ID:       string(i),
						Metadata: metadata,
						Content:  content,
					})
				if err != nil {
					log.Println("Fail !!!!!")
					panic(err)
				}
			}
		}
		time.Sleep(700 * time.Microsecond)
		counter = 0

	}
}
