package main

import (
	"fmt"
	"gocoder/input"
	"log"
)

func main() {
	context := input.NewContext()
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

	c := context.ReadInput()

	for elem := range c {
		fmt.Printf("Packet: %d\n", elem.Size())
	}
}