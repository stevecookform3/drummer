package audio

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"log"
	"os"
)

type Sample struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeekCloser
}

type FileType int

const (
	WAV FileType = iota
	MP3
)

// Load an sample, in either WAV or MP3 format
func NewSample(filePath string, fileType FileType) *Sample {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Couldnt open file %s", err)
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format

	if fileType == MP3 {
		streamer, format, err = mp3.Decode(f)
	} else if fileType == WAV {
		streamer, format, err = wav.Decode(f)
	}

	if err != nil {
		log.Fatalf("Couldnt decode file %s", err)
	}

	return &Sample{format.SampleRate, streamer}
}

func (s *Sample) Play(volume float64, pitch float64) {
	s.streamer.Seek(0)
	resampler := beep.ResampleRatio(4, pitch, s.streamer)
	volumeEffect := &effects.Volume{Streamer: resampler, Base: 2, Volume: volume}
	speaker.Play(volumeEffect)
}

func (s *Sample) Close() {
	s.streamer.Close()
}
