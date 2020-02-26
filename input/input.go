package input

import (
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"unsafe"
)

func init() {
	avformat.AvRegisterAll()
}

type Context struct {
	FormatCtx      *avformat.Context
	Filename       string
	DecodeContexts []*avcodec.Context
}

func NewContext() Context {
	return Context{
		FormatCtx: avformat.AvformatAllocContext(),
		Filename:  "",
		DecodeContexts: nil,
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
			if err == avutil.AvErrorEOF {
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

func (c *Context) DecodeStream(stream <-chan *avcodec.Packet) {
	frame := avutil.AvFrameAlloc()

	for packet := range stream {
		sendErr := c.sendToDecoder(packet)
		index := packet.StreamIndex()

		for sendErr >= 0 {
			resp := c.getFromDecoder(index, frame)

			if resp == avutil.AvErrorEAGAIN || resp == avutil.AvErrorEOF {
				//fmt.Printf("Resp: %s\n", avutil.ErrorFromCode(resp))
				break
			} else if resp < 0 {
				fmt.Printf("Error getting frame from decoder: %s\n",
					avutil.ErrorFromCode(resp))
			}
		}

		//width, height, _, _ := avutil.AvFrameGetInfo(frame)
		//fmt.Printf("Frame width/height: %dx%d\n", width, height)
	}
}

func (c *Context) OpenInput(path string) error {
	c.Filename = path
	err := avformat.AvformatOpenInput(&c.FormatCtx, c.Filename, nil, nil)

	if err != 0 {
		return fmt.Errorf("OpenInput: %v", avutil.ErrorFromCode(err))
	}

	for _, stream := range c.FormatCtx.Streams() {
		codecContext := stream.Codec()

		codec := avcodec.AvcodecFindDecoder(avcodec.CodecId(codecContext.GetCodecId()))

		if codec == nil {
			return fmt.Errorf("could not find decoder for %v", codecContext.GetCodecId())
		}

		codecName := avcodec.AvcodecGetName(avcodec.CodecId(codecContext.GetCodecId()))
		fmt.Printf("Opening decoder: %s\n", codecName)

		decodeContext := codec.AvcodecAllocContext3()
		err = decodeContext.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(codecContext)))

		if err != 0 {
			return fmt.Errorf("OpenInput: %v", avutil.ErrorFromCode(err))
		}

		err = decodeContext.AvcodecOpen2(codec, nil)

		if err < 0 {
			return fmt.Errorf("OpenInput: %v", avutil.ErrorFromCode(err))
		}

		c.DecodeContexts = append(c.DecodeContexts, decodeContext)
	}

	return nil
}

func (c *Context) getFromDecoder(index int, frame *avutil.Frame) int {
	decodingContext := c.DecodeContexts[index]
	return decodingContext.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
}

func (c *Context) sendToDecoder(packet *avcodec.Packet) int {
	streamIndex := packet.StreamIndex()
	decodingContext := c.DecodeContexts[streamIndex]

	err := decodingContext.AvcodecSendPacket(packet)

	if err < 0 {
		fmt.Printf("Error sending frame to decoder: %s\n",
			avutil.ErrorFromCode(err))
	}

	return err
}

func (c *Context) CloseInput() {
	c.FormatCtx.AvformatCloseInput()
}

func (c *Context) FindStreamInfo() {
	c.FormatCtx.AvformatFindStreamInfo(nil)
}

func (c *Context) DumpStreamInfo() {
	c.FormatCtx.AvDumpFormat(0, c.Filename, 0)
}