package input

import (
	"testing"
)

const (
	filename = "../assets/small.mp4"
)

func TestOpenClose(t *testing.T) {
	context := NewContext()

	if context.FormatCtx == nil {
		t.Error("Format context not created")
	}

	t.Run("OpenFileAndCodecs", func(t *testing.T) {
		err := context.OpenInput(filename)

		if err != nil {
			t.Error(err)
		}

		nrStreams := len(context.Streams)
		if nrStreams != 2 {
			t.Errorf("Missing streams, found %d\n", nrStreams)
		}
	})

	t.Run("DecodeInput", func(t *testing.T) {
		c := context.ReadInput()

		frames := context.DecodeStream(c)
		count := 0
		for frame := range frames {
			width := frame.Width()
			height := frame.Height()

			if width != 0 && height != 0 {
				count++
			}
		}

		if count != 166 {
			t.Errorf("Incorrect number %d of frames decoded\n", count)
		}
	})

	t.Run("CloseInput", func(t *testing.T) {
		context.CloseInput()
	})
}
