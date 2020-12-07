package audio

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
)

func NewOutput() {
	playbackFormat := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}
	speaker.Init(playbackFormat.SampleRate, playbackFormat.SampleRate.N(time.Second/30))
}

func Close() {
	speaker.Close()
}
