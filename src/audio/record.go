package audio

import (
	"code.google.com/p/portaudio-go/portaudio"
	"encoding/binary"
	_ "fmt"
	"os"
	"time"
)

func init() {
	portaudio.Initialize()
}

const (
	DefaultSampleRate = 44100
)

type RecordOptions struct {
	FilePath         string
	MaxDuration      time.Duration
	StopSignal       chan int
	Callback         func(nsamples int)
	CallbackInterval int     // default: 4410
	SampleRate       float64 // default: 44100
	InputChannels    int     // number of input channels. default: 1
}

// Record audio into an AIFF file using PortAudio
// Base on example code from: https://code.google.com/p/portaudio-go/source/browse/portaudio/examples/record.go
// func RecordAIFF(p string, dur time.Duration, sig chan int) error {
func RecordAIFF(options RecordOptions) error {
	if options.SampleRate <= 0 {
		options.SampleRate = DefaultSampleRate
	}
	if options.InputChannels <= 0 {
		options.InputChannels = 1
	}
	if options.CallbackInterval <= 0 {
		options.CallbackInterval = int(options.SampleRate / 10)
	}

	f, err := os.Create(options.FilePath)
	if err != nil {
		return err
	}

	// form chunk
	if _, err = f.WriteString("FORM"); err != nil {
		return err
	}

	err = binary.Write(f, binary.BigEndian, int32(0)) //total bytes
	if err != nil {
		return err
	}
	if _, err = f.WriteString("AIFF"); err != nil {
		return err
	}

	// common chunk
	if _, err = f.WriteString("COMM"); err != nil {
		return err
	}
	binary.Write(f, binary.BigEndian, int32(18))                       //size
	binary.Write(f, binary.BigEndian, int16(1))                        //channels
	binary.Write(f, binary.BigEndian, int32(0))                        //number of samples
	binary.Write(f, binary.BigEndian, int16(32))                       //bits per sample
	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	if err != nil {
		return err
	}

	// sound chunk
	if _, err = f.WriteString("SSND"); err != nil {
		return err
	}
	binary.Write(f, binary.BigEndian, int32(0)) //size
	binary.Write(f, binary.BigEndian, int32(0)) //offset
	binary.Write(f, binary.BigEndian, int32(0)) //block
	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		if _, err = f.Seek(4, 0); err != nil {
			panic(err)
		}
		if err := binary.Write(f, binary.BigEndian, int32(totalBytes)); err != nil {
			panic(err)
		}
		if _, err = f.Seek(22, 0); err != nil {
			panic(err)
		}
		if err := binary.Write(f, binary.BigEndian, int32(nSamples)); err != nil {
			panic(err)
		}
		if _, err = f.Seek(42, 0); err != nil {
			panic(err)
		}
		if err := binary.Write(f, binary.BigEndian, int32(4*nSamples+8)); err != nil {
			panic(err)
		}
		f.Close()
	}()

	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(options.InputChannels, 0, options.SampleRate, len(in), in)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return err
	}

	cbSamples := 0 // sample count of last callback

	for {
		if err := stream.Read(); err != nil {
			return err
		}
		err = binary.Write(f, binary.BigEndian, in)
		if err != nil {
			return err
		}
		nSamples += len(in)

		if options.StopSignal != nil {
			select {
			case <-options.StopSignal:
				return nil
			default:
			}
		}
		if options.Callback != nil {
			if cbSamples == 0 || nSamples-cbSamples > options.CallbackInterval {
				options.Callback(nSamples)
				cbSamples = nSamples
			}
		}
		if options.MaxDuration > 0 {
			if float64(nSamples) > float64(options.MaxDuration/time.Second)*float64(options.SampleRate) {
				break
			}
		}
	}

	if err = stream.Stop(); err != nil {
		return err
	}

	return nil
}
