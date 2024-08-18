package main

import (
	vsegpt "github.com/saintbyte/vsegpt_api"
	"log"
)

func main() {
	vg := vsegpt.NewVseGpt()
	s, err := vg.Ask("Сколько в море рыбы?")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)
}
