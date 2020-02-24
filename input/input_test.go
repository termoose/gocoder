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

		nrStreams := len(context.DecodeContexts)
		if nrStreams != 2 {
			t.Errorf("Missing streams, found %d\n", nrStreams)
		}
	})

	t.Run("ReadInput", func(t *testing.T) {
		c := context.ReadInput()

		size := 0
		for elem := range c {
			size += elem.Size()
		}

		// Check that we have demuxed everything
		if size != 379872 {
			t.Errorf("File size incorrect: %d bytes", size)
		}
	})

	t.Run("CloseInput", func(t *testing.T) {
		context.CloseInput()
	})
}
