package input

import (
	"fmt"
	"github.com/asticode/goav/avcodec"
	"github.com/asticode/goav/avformat"
	"github.com/asticode/goav/avutil"
	"gocoder/encode"
)

type Stream struct {
	DecodeContext *avcodec.Context
	Params        *avcodec.CodecParameters
}

type Context struct {
	FormatCtx      *avformat.Context
	Filename       string
	Streams        []Stream
}

func NewContext() Context {
	return Context{
		FormatCtx: avformat.AvformatAllocContext(),
		Filename:  "",
		Streams: nil,
	}
}

func (c *Context) ReadInput() chan *avcodec.Packet {
	outBuffer := make(chan *avcodec.Packet, 50)

	go func() {
		for {
			// FIXME: deref this?
			packet := avcodec.AvPacketAlloc()
			err := c.FormatCtx.AvReadFrame(packet)

			// FIXME: find a good way to communicate this out through the channel?
			if err == avutil.AVERROR_EOF {
				// Send the EOF packet to signal EOF to decoder thread
				outBuffer <- packet
				close(outBuffer)
				return
			} else if err < 0 {
				close(outBuffer)
				return
			}

			outBuffer <- packet
		}
	}()

	return outBuffer
}

func (c *Context) DecodeStream(stream <-chan *avcodec.Packet) chan encode.Frame {
	outBuffer := make(chan encode.Frame, 50)

	go func() {
		defer close(outBuffer)

		for packet := range stream {
			_ = c.sendToDecoder(packet)
			index := packet.StreamIndex()

			for err := 0; err >= 0; {
				frame := avutil.AvFrameAlloc()
				err = c.getFromDecoder(index, frame)

				if err == avutil.AVERROR_EAGAIN {
					break
				} else if err == avutil.AVERROR_EOF {
					// Send EOF frame to signal EOF to encoders?
					//outBuffer <- frame
					return
				} else if err < 0 {
					fmt.Printf("Error getting frame from decoder: %s\n",
						avutil.AvStrerr(err))
					return
				}

				if index == 0 {
					outBuffer <- encode.NewFrame(frame, c.mediaType(index))
				}
			}
		}
	}()

	return outBuffer
}

func (c *Context) OpenInput(path string) error {
	c.Filename = path
	err := avformat.AvformatOpenInput(&c.FormatCtx, c.Filename, nil, nil)

	if err != 0 {
		return fmt.Errorf("OpenInput: %v", avutil.AvStrerr(err))
	}

	for _, stream := range c.FormatCtx.Streams() {
		params := stream.CodecParameters()
		codec := avcodec.AvcodecFindDecoder(params.CodecId())

		if codec == nil {
			return fmt.Errorf("could not find decoder for %v", params.CodecId())
		}

		decodeContext := codec.AvcodecAllocContext3()
		err = avcodec.AvcodecParametersToContext(decodeContext, params)

		if err != 0 {
			return fmt.Errorf("OpenInput: %v", avutil.AvStrerr(err))
		}

		err = decodeContext.AvcodecOpen2(codec, nil)

		if err < 0 {
			return fmt.Errorf("OpenInput: %v", avutil.AvStrerr(err))
		}

		c.Streams = append(c.Streams, Stream{
			DecodeContext: decodeContext,
			Params: params,
		})
	}

	return nil
}

func (c *Context) mediaType(index int) avformat.MediaType {
	return avformat.MediaType(c.Streams[index].Params.CodecType())
}

func (c *Context) getFromDecoder(index int, frame *avutil.Frame) int {
	stream := c.Streams[index]
	return avcodec.AvcodecReceiveFrame(stream.DecodeContext, frame)
}

func (c *Context) sendToDecoder(packet *avcodec.Packet) int {
	streamIndex := packet.StreamIndex()
	stream := c.Streams[streamIndex]

	err := avcodec.AvcodecSendPacket(stream.DecodeContext, packet)
	if err < 0 {
		fmt.Printf("Error sending frame to decoder: %s\n",
			avutil.AvStrerr(err))
	}

	return err
}

func (c *Context) CloseInput() {
	avformat.AvformatCloseInput(c.FormatCtx)
}

func (c *Context) FindStreamInfo() {
	c.FormatCtx.AvformatFindStreamInfo(nil)
}

func (c *Context) DumpStreamInfo() {
	c.FormatCtx.AvDumpFormat(0, c.Filename, 0)
}