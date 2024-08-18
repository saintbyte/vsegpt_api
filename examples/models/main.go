package main

import (
	vsegpt "github.com/saintbyte/vsegpt_api"
	"log"
)

func main() {
	vg := vsegpt.NewVseGpt()
	models, err := vg.GetModels()
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		log.Println(model.ID)
		log.Println(model.Name)
		log.Println(model.ContextLength)
		log.Println(model.ModelPricing)
		log.Println("----------------")
	}
}
