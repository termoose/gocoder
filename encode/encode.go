package encode

import (
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"unsafe"
)

type Encoder interface {
	Encode(frame *avutil.Frame) *avcodec.Packet
}

type Video struct {
	Codec   *avcodec.Codec
	Context *avcodec.Context
}

type Audio struct {
	Codec   *avcodec.Codec
	Context *avcodec.Context
}

func (v *Video) Encode(stream <-chan *avutil.Frame) chan *avcodec.Packet {
	outBuffer := make(chan *avcodec.Packet, 150)

	go func() {
		for frame := range stream {
			codecFrame := (*avcodec.Frame)(unsafe.Pointer(frame))
			_ = v.Context.AvcodecSendFrame(codecFrame)

			for err := 0; err >= 0; {
				packet := avcodec.AvPacketAlloc()
				err = v.Context.AvcodecReceivePacket(packet)

				if err == avutil.AvErrorEAGAIN {
					break
				} else if err == avutil.AvErrorEOF {
					fmt.Println("EOF encode")
					close(outBuffer)
					return
				} else if err < 0 {
					fmt.Printf("Error getting frame from encoder: %s\n",
					avutil.ErrorFromCode(err))
					close(outBuffer)
					return
				}

				outBuffer <- packet
			}
		}
	}()

	return outBuffer
}

func (v *Video) SetOptions(width, height int) {
	v.Context.SetEncodeParams2(width, height, avcodec.AV_PIX_FMT_YUV,
		true, 25)

	// Hack for setting bitrate? Remove this in private fork
	ctx := (*avformat.CodecContext)(unsafe.Pointer(v.Context))
	ctx.SetBitRate(1000000)
	ctx.SetTimeBase(avcodec.NewRational(1, 25))

	v.Context.AvcodecOpen2(v.Codec, nil)
}

func NewVideoEncoder() Video {
	codec := avcodec.AvcodecFindEncoderByName("libx264")

	return Video{
		Context: codec.AvcodecAllocContext3(),
	}

	//codec := avcodec.AvcodecFindEncoderByName("libx264")
	//
	//v.Context := codec.AvcodecAllocContext3()
	//context.SetEncodeParams2(800, 600, avcodec.AV_PIX_FMT_YUV, true, 25)
	//
	//err := context.AvcodecOpen2(codec, nil)
	//
	//if err < 0 {
	//	return fmt.Errorf("NewEncoder: %w, ", avutil.ErrorFromCode(err))
	//}
	//
	//return nil
}

func (v *Video) EncodeStream(stream <-chan *avutil.Frame) {
	for _ = range stream {

	}
}