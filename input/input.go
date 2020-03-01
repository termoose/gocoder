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

func (c *Context) DecodeStream(stream <-chan *avcodec.Packet) chan *avutil.Frame {
	outBuffer := make(chan *avutil.Frame, 50)

	go func() {
		for packet := range stream {
			_ = c.sendToDecoder(packet)
			index := packet.StreamIndex()

			for err := 0; err >= 0; {
				frame := avutil.AvFrameAlloc()
				err = c.getFromDecoder(index, frame)

				if err == avutil.AvErrorEAGAIN {
					break
				} else if err == avutil.AvErrorEOF {
					fmt.Printf("EOF decode")
					close(outBuffer)
					return
				} else if err < 0 {
					fmt.Printf("Error getting frame from decoder: %s\n",
						avutil.ErrorFromCode(err))
					close(outBuffer)
					return
				}

				outBuffer <- frame
			}
		}
	}()

	return outBuffer
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