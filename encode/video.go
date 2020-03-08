package encode

import (
	"fmt"
	"github.com/asticode/goav/avcodec"
	"github.com/asticode/goav/avutil"
)

type Video struct {
	Codec   *avcodec.Codec
	Context *avcodec.Context
}

func NewVideoEncoder() Video {
	codec := avcodec.AvcodecFindEncoderByName("libx264")

	return Video{
		Context: codec.AvcodecAllocContext3(),
	}
}

func (v *Video) Encode(stream <-chan *avutil.Frame) chan *avcodec.Packet {
	outBuffer := make(chan *avcodec.Packet, 50)

	go func() {
		defer close(outBuffer)

		for frame := range stream {
			// Reset all frame types to avoid weird GOP's
			frame.SetPictType(avutil.AV_PICTURE_TYPE_NONE)

			ret := avcodec.AvcodecSendFrame(v.Context, frame)

			if ret < 0 {
				fmt.Printf("Error sending frame to encoder: %s\n",
					avutil.AvStrerr(ret))
				return
			}

			for err := 0; err >= 0; {
				packet := avcodec.AvPacketAlloc()
				err = avcodec.AvcodecReceivePacket(v.Context, packet)

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
	v.Context.SetBitRate(1000000)
	v.Context.SetTimeBase(avutil.NewRational(1, 25))
	v.Context.SetFramerate(avutil.NewRational(25, 1))
	v.Context.SetPixFmt(avutil.AV_PIX_FMT_YUV420P)
	v.Context.SetGopSize(25)
	v.Context.SetWidth(width)
	v.Context.SetHeight(height)

	err := v.Context.AvcodecOpen2(v.Codec, nil)

	if err < 0 {
		fmt.Printf("Error opening codec: %v\n", avutil.AvStrerr(err))
	}
}