package encode

import (
	"fmt"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
)

type EncodingContext struct {
	FormatCtx *avformat.Context
	Filename  string
}

func NewEncodingContext() EncodingContext {
	return EncodingContext{
		FormatCtx: avformat.AvformatAllocContext(),
		Filename:  "",
	}
}

func (c *EncodingContext) OpenInput(path string) error {
	c.Filename = path
	err := avformat.AvformatOpenInput(&c.FormatCtx, c.Filename, nil, nil)

	if err != 0 {
		return fmt.Errorf("OpenInput: %v", avutil.ErrorFromCode(err))
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