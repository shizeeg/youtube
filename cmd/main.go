package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shizeeg/youtube"
)

func main() {
	var id string
	if len(os.Args) > 0 {
		id = os.Args[1]
	}
	key := os.Getenv("YOUTUBEDATAKEY")
	if key == "" {
		log.Fatalln("YOUTUBEDATAKEY environment variable not set")
	}
	if len(id) != 11 {
		log.Fatalf("%q - wrong YouTube video id (must be 11 characters)", id)
	}
	duration, err := youtube.GetDuration(key, id)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(duration)
}
