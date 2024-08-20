package main

import (
	vsegpt "github.com/saintbyte/vsegpt_api"
	"log"
)

func main() {
	vg := vsegpt.NewVseGpt()
	vg.Temperature = 1.2
	s, err := vg.Ask("Сколько в море рыбы?")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)
}
