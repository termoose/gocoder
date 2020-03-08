package encode

import (
	"github.com/asticode/goav/avcodec"
	"github.com/asticode/goav/avutil"
)

type FrameProcessor interface {
	Encode(frame *avutil.Frame) *avcodec.Packet
}

type EncoderContext struct {

}

type Audio struct {
	Codec   *avcodec.Codec
	Context *avcodec.Context
}

func (v *EncoderContext) Process(stream <-chan *avutil.Frame) {
	for _ = range stream {

	}
}
