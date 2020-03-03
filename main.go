package main

import (
	"fmt"
	"github.com/giorgisio/goav/avutil"
	"gocoder/input"
	"gocoder/encode"
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

	// Demux
	c := context.ReadInput()

	// Decode
	decodedFrames := context.DecodeStream(c)
	for elem := range decodedFrames {
		width, height, _, _ := avutil.AvFrameGetInfo(elem)
		fmt.Printf("Frame %dx%d\n", width, height)
	}

	video := encode.NewVideoEncoder()
	video.SetOptions(800, 600)
}