package main

import (
	"github.com/giorgisio/goav/avformat"
	"gocoder/encode"
	"log"
)

func main() {
	avformat.AvRegisterAll()
	context := encode.NewEncodingContext()
	filename := "small.mp4"

	err := context.OpenInput(filename)
	defer context.CloseInput()

	if err != nil {
		log.Printf("Could not open file: %v\n", err)
		return
	}

	context.FindStreamInfo()
	context.DumpStreamInfo()
}