package encode

import (
	"testing"
)

const (
	filename = "../assets/small.mp4"
)

func TestOpenClose(t *testing.T) {
	context := NewEncodingContext()

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

	t.Run("CloseInput", func(t *testing.T) {
		context.CloseInput()
	})
}
