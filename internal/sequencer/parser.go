package sequencer

import (
	"bufio"
	"github.com/go-yaml/yaml"
	"strings"
)

type Config struct {
	Instruments map[string]Instrument
	Sequence    Sequence
}

type Instrument struct {
	Sample string
	Volume float64
	Pitch  float64
}

type Sequence struct {
	Tempo  int
	SubDiv int
	Seq    string
}

func ParseConfig(yamlConfig []byte) (Config, error) {
	config := Config{}
	err := yaml.Unmarshal(yamlConfig, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

type Beat struct {
	Instrument []string
}

func ParseSequence(seq string) []Beat {

	var sequence []Beat
	scanner := bufio.NewScanner(strings.NewReader(seq))
	for scanner.Scan() {
		instruments := strings.Split(scanner.Text(), "|")
		if sequence == nil {
			sequence = make([]Beat, len(instruments)-2) //ignore empty leading & trailing beats
		}
		for index, instrument := range instruments {
			if index == 0 || index == len(instruments)-1 || strings.TrimSpace(instrument) == "" {
				continue
			}
			sequence[index-1].Instrument = append(sequence[index-1].Instrument, strings.TrimSpace(instrument))
		}
	}
	return sequence
}
