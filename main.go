package main

import (
	"fmt"
	"gocoder/encode"
	"log"
)

func main() {
	context := encode.NewEncodingContext()
	filename := "assets/small.mp4"

	err := context.OpenInput(filename)
	defer context.CloseInput()

	if err != nil {
		log.Printf("Could not open file: %v\n", err)
		return
	}

	context.FindStreamInfo()
	context.DumpStreamInfo()

	for i, ctx := range context.DecodeContexts {
		fmt.Printf("Stream %d codec %p\n", i, ctx)
	}
}