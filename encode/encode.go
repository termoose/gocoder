package encode

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

type EncodingContext struct {
	FormatCtx      *avformat.Context
	Filename       string
	DecodeContexts []*avcodec.Context
}

func NewEncodingContext() EncodingContext {
	return EncodingContext{
		FormatCtx: avformat.AvformatAllocContext(),
		Filename:  "",
		DecodeContexts: nil,
	}
}

func (c *EncodingContext) OpenInput(path string) error {
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

func (c *EncodingContext) CloseInput() {
	c.FormatCtx.AvformatCloseInput()
}

func (c *EncodingContext) FindStreamInfo() {
	c.FormatCtx.AvformatFindStreamInfo(nil)
}

func (c *EncodingContext) DumpStreamInfo() {
	c.FormatCtx.AvDumpFormat(0, c.Filename, 0)
}