package main

import (
	"fmt"
	"github.com/form3tech/drummer/internal/audio"
	"github.com/form3tech/drummer/internal/sequencer"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage: drummer example.yml")
		return
	}

	configYaml, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Couldnt open file %s", err)
	}

	config, err := sequencer.ParseConfig(configYaml)
	if err != nil {
		log.Fatalf("Couldnt parse sequence %s", err)
	}

	audio.NewOutput()
	defer audio.Close()

	//load instruments
	samples := make(map[string]*audio.Sample)
	for key, instrument := range config.Instruments {
		samples[key] = audio.NewSample(instrument.Sample, audio.WAV)
	}

	//load sequence
	beats := sequencer.ParseSequence(config.Sequence.Seq)
	fmt.Println("Drum Machine Sequence:")
	fmt.Println(config.Sequence.Seq)

	//play sequence
	beatInterval := int(time.Minute) / (config.Sequence.Tempo * config.Sequence.SubDiv)
	ticker := time.NewTicker(time.Duration(beatInterval))
	done := make(chan bool)

	index := 0
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, instrumentKey := range beats[index].Instrument {
					samples[instrumentKey].Play(config.Instruments[instrumentKey].Volume, config.Instruments[instrumentKey].Pitch)
				}
				index++
				index = index % len(beats)
			}
		}
	}()

	fmt.Println("Press Enter to exit! ¯\\_(ツ)_/¯")
	var input string
	fmt.Scanln(&input)

	ticker.Stop()
	done <- true

	for _, sample := range samples {
		sample.Close()
	}
}
