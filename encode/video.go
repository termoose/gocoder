package encode

import (
	"fmt"
	"github.com/asticode/goav/avcodec"
	"github.com/asticode/goav/avutil"
)

type Video struct {
	codec   *avcodec.Codec
	context *avcodec.Context
}

func NewVideoEncoder() Video {
	codec := avcodec.AvcodecFindEncoderByName("libx264")

	return Video{
		codec:   codec,
		context: codec.AvcodecAllocContext3(),
	}
}

func (v *Video) Encode(stream <-chan Frame) chan *avcodec.Packet {
	outBuffer := make(chan *avcodec.Packet, 50)

	go func() {
		defer close(outBuffer)

		for frame := range stream {
			avFrame := frame.AVFrame
			// Reset all frame types to avoid weird GOP's
			avFrame.SetPictType(avutil.AV_PICTURE_TYPE_NONE)

			ret := avcodec.AvcodecSendFrame(v.context, avFrame)

			if ret < 0 {
				fmt.Printf("Error sending frame to encoder: %s\n",
					avutil.AvStrerr(ret))
				return
			}

			for err := 0; err >= 0; {
				packet := avcodec.AvPacketAlloc()
				err = avcodec.AvcodecReceivePacket(v.context, packet)

				if err == avutil.AVERROR_EAGAIN {
					break
				} else if err == avutil.AVERROR_EOF {
					fmt.Println("EOF encode")
					return
				} else if err < 0 {
					fmt.Printf("Error getting frame from encoder: %s\n",
						avutil.AvStrerr(err))
					return
				}

				outBuffer <- packet
			}
		}
	}()

	return outBuffer
}

func (v *Video) SetOptions(width, height int) {
	v.context.SetBitRate(1000000)
	v.context.SetTimeBase(avutil.NewRational(1, 25))
	v.context.SetFramerate(avutil.NewRational(25, 1))
	v.context.SetPixFmt(avutil.AV_PIX_FMT_YUV420P)
	v.context.SetGopSize(25)
	v.context.SetWidth(width)
	v.context.SetHeight(height)

	err := v.context.AvcodecOpen2(v.codec, nil)

	if err < 0 {
		fmt.Printf("Error opening codec: %v\n", avutil.AvStrerr(err))
	}
}