package main

import (
	"fmt"
	_ "github.com/giorgisio/goav/avutil"
	"gocoder/encode"
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

	for i, ctx := range context.Streams {
		fmt.Printf("Stream %d stream %p\n", i, ctx)
	}

	// Demux
	c := context.ReadInput()

	// Decode
	decodedFrames := context.DecodeStream(c)
	//for elem := range decodedFrames {
	//	width, height, _, _ := avutil.AvFrameGetInfo(elem)
	//	fmt.Printf("Frame %dx%d\n", width, height)
	//}

	video := encode.NewVideoEncoder()
	video.SetOptions(560, 320)
	encodedFrames := video.Encode(decodedFrames)

	for elem := range encodedFrames {
		fmt.Printf("size: %d\n", elem.Size())
	}
}