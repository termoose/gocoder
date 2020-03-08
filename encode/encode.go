package encode

import (
	"github.com/asticode/goav/avcodec"
	"github.com/asticode/goav/avformat"
	"github.com/asticode/goav/avutil"
)

type Frame struct {
	avFrame   *avutil.Frame
	frameType avformat.MediaType
}

type FrameProcessor interface {
	Encode(frame Frame) *avcodec.Packet
}

type EncoderContext struct {

}

type Audio struct {
	Codec   *avcodec.Codec
	Context *avcodec.Context
}

func NewFrame(avFrame *avutil.Frame, frameType avformat.MediaType) Frame {
	return Frame{
		avFrame: avFrame,
		frameType: frameType,
	}
}

func (v *EncoderContext) Process(stream <-chan Frame) {
	for frame := range stream {
		switch frame.frameType {
		case avutil.AVMEDIA_TYPE_VIDEO:

		}

	}
}
