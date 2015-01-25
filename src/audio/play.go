package audio

import (
	"code.google.com/p/portaudio-go/portaudio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

func init() {
	portaudio.Initialize()
}

var (
	ErrBadFileFormat = errors.New("bad file format")
)

// Play an AIFF file using PortAudio
// Base on example code from: https://code.google.com/p/portaudio-go/source/browse/portaudio/examples/play.go
func PlayAIFF(p string, sig chan int) error {

	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	id, data, err := readChunk(f)
	if err != nil {
		return err
	}
	if id.String() != "FORM" {
		return ErrBadFileFormat
	}
	_, err = data.Read(id[:])
	if err != nil {
		return err
	}
	if id.String() != "AIFF" {
		return ErrBadFileFormat
	}

	var c commonChunk
	var audio io.Reader
	for {
		id, chunk, err := readChunk(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch id.String() {
		case "COMM":
			err = binary.Read(chunk, binary.BigEndian, &c)
			if err != nil {
				return err
			}
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		default:
			fmt.Printf("ignoring unknown chunk '%s'\n", id)
		}
	}

	//assume 44100 sample rate, mono, 32 bit

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	if err != nil {
		return err
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		return err
	}

	defer stream.Stop()
	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Write()
		if err != nil {
			return err
		}

		select {
		case <-sig:
			return nil
		default:
		}
	}

	return nil
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
